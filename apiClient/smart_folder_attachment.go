package apiClient

import (
	"fmt"

	errorsHandler "github.com/turbot/terraform-provider-turbot/errors"
	"sync"
	"time"
)

// The attachSmartFolders and detachSmartFolders mutations are read-modify-write on the TARGET
// resource's attachment list: the server reads the current list, applies the change and writes
// it back. Concurrent calls for the same target therefore lose updates.
//
// Measured against a live workspace, attaching six policy packs to one folder:
//
//	terraform apply                    -> 4 of 6 attached, "Apply complete", NO error
//	terraform apply -parallelism=1     -> 6 of 6 attached, no error
//
// The losing attachments are reported as created, so Terraform records them in state and the
// drift only surfaces on a later plan. Serialising writes per target removes the race. Attaching
// many packs to one resource in a single apply becomes serial, which is slower but correct; the
// lock is per target, so unrelated targets still attach concurrently.
var attachmentTargetLocks sync.Map // target id -> *sync.Mutex

// lockAttachmentTarget serialises attachment writes for a single target. It returns the unlock
// func so callers can defer it.
func lockAttachmentTarget(target string) func() {
	value, _ := attachmentTargetLocks.LoadOrStore(target, &sync.Mutex{})
	mutex := value.(*sync.Mutex)
	mutex.Lock()
	return mutex.Unlock
}

// attachmentTarget pulls the target identifier out of a mutation input. Returns "" when absent,
// in which case the caller skips locking rather than serialising every attachment in the process
// behind one key.
func attachmentTarget(input map[string]interface{}) string {
	target, ok := input["resource"].(string)
	if !ok {
		return ""
	}
	return target
}

// attachmentLockKey resolves target to its numeric Turbot id, so a target written as an id in one
// attachment and as an aka in another takes the SAME lock. Keying on the raw string instead would
// let the race survive any config that mixes the two forms - an ordinary way for a config to drift
// over time.
//
// This costs one extra read per attachment write. verifyAttachmentState reads attachedSmartFolders,
// which does not yield the target's own id, so it cannot be reused here; and the key must be known
// BEFORE the mutation, not after. Resolution uses `resource(id:)`, which the caller is necessarily
// entitled to since it holds permissions on the attachment target. On failure it falls back to the
// raw value: a weaker lock, but a read error here must not block an attach the caller may make.
// Resolved lock keys are cached: attaching N packs to one target would otherwise resolve the same
// target N times, and because the key is needed before the lock those reads all fire concurrently
// on the first pass. Caching also makes the fallback STICKY, which matters for correctness: without
// it a transient read failure on one goroutine gives that writer the raw string while its siblings
// get the numeric id, so they take different locks and the race is live for exactly that window.
// Caching the first answer makes every writer NAMING THE TARGET THE SAME WAY agree. Across forms it
// can still diverge - if resolution fails for an aka while succeeding for the id, those two land on
// different keys - which needs a transient read failure and a config mixing both forms and
// concurrent writes to one target, and is backstopped by verification below.
//
// Safe to cache for the life of the process: a resource's numeric id never changes, and the
// provider process is short-lived, so staleness has no window in which to matter.
var attachmentLockKeys sync.Map // raw target identifier -> resolved numeric id

func (client *Client) attachmentLockKey(target string) string {
	if cached, ok := attachmentLockKeys.Load(target); ok {
		return cached.(string)
	}
	key := target
	if resource, err := client.ReadResource(target, nil); err == nil && resource.Turbot.Id != "" {
		key = resource.Turbot.Id
	}
	actual, _ := attachmentLockKeys.LoadOrStore(target, key)
	return actual.(string)
}

// Post-write confirmation bounds. The attachment list is not immediately consistent after a write,
// so a single read can report a successful write as not yet applied. The delay ramps
// 250ms -> 500ms -> 1s -> 2s (3.75s total) rather than sitting at a flat 2s: a retry is the
// expected path rather than the exceptional one, and the wait is paid while holding the target's
// lock, so a coarse first backoff would queue every same-target write behind it.
// Vars rather than consts purely so tests can shrink the budget; nothing in the provider reassigns
// them.
var (
	verifyAttachmentAttempts  = 5
	verifyAttachmentBaseDelay = 250 * time.Millisecond
)

// verifyBackoff returns the wait before retry number attempt (1-based). Attempts below 1 return
// zero rather than panicking on a negative shift - unreachable through the current call site, but a
// panic in a provider is one refactor away and costs nothing to make structurally impossible.
func verifyBackoff(attempt int) time.Duration {
	if attempt < 1 {
		return 0
	}
	return verifyAttachmentBaseDelay << (attempt - 1)
}

func (client *Client) CreateSmartFolderAttachment(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateSmartFolderAttachResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	if target := attachmentTarget(input); target != "" {
		defer lockAttachmentTarget(client.attachmentLockKey(target))()
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		// handleCreateError's not-found branch reports input["parent"], which an attachment input
		// does not carry - it holds only `resource` and `smartFolders` - so it renders
		// "parent resource not found: %!s(<nil>)". Name the things that can actually be missing.
		if errorsHandler.NotFoundError(err) {
			return nil, fmt.Errorf("error creating smart folder attachment: resource or policy pack not found (resource: %v, policy packs: %v)", input["resource"], input["smartFolders"])
		}
		return nil, client.handleCreateError(err, input, "smart folder attachment")
	}

	// Locking prevents the race within one process, but nothing serialises separate terraform runs
	// against the same target. Confirm the attachment landed so a lost write fails loudly here
	// rather than being recorded in state and reappearing as drift on a later plan.
	if err := client.verifyAttachment(input); err != nil {
		return nil, err
	}

	return &responseData.SmartFolderAttach.Turbot, nil
}

func (client *Client) DeleteSmartFolderAttachment(input map[string]interface{}) error {
	query := detachSmartFolderAttachment()
	var responseData interface{}

	variables := map[string]interface{}{
		"input": input,
	}

	// Detach is the same read-modify-write on the target's list, so it takes the same lock.
	if target := attachmentTarget(input); target != "" {
		defer lockAttachmentTarget(client.attachmentLockKey(target))()
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting smart folder attachment: %s", err.Error())
	}

	// A lost detach is the more dangerous direction, and the only one with no backstop. A lost
	// attach is caught by verification here, or failing that by drift on the next plan, which
	// re-attaches. A lost detach is not: the mutation reports success, Terraform drops the resource
	// from state, and Exists is never consulted again for something that left state - so the pack
	// stays attached and keeps evaluating its policies against the target indefinitely.
	return client.verifyDetachment(input)
}

// verifyAttachment confirms every pack named in input IS attached to the target.
func (client *Client) verifyAttachment(input map[string]interface{}) error {
	return client.verifyAttachmentState(input, true)
}

// verifyDetachment confirms none of the packs named in input REMAIN attached to the target.
func (client *Client) verifyDetachment(input map[string]interface{}) error {
	return client.verifyAttachmentState(input, false)
}

// verifyAttachmentState polls the target's attachment list until every pack in input reaches
// wantAttached, or the attempts run out. A nil return means confirmed, or that verification could
// not be performed at all - a mutation the server accepted is not failed merely because the check
// could not run.
//
// Truncation inverts meaning between the two directions: verifying an attach, absence under
// truncation proves nothing, because the pack may sit on a page that was not returned; verifying a
// detach, presence is the positive signal and truncation can only hide a pack that is still
// attached. Bailing out is the conservative choice in both directions, so both take it.
func (client *Client) verifyAttachmentState(input map[string]interface{}, wantAttached bool) error {
	target := attachmentTarget(input)
	packs := attachmentPacks(input)
	if target == "" || len(packs) == 0 {
		return nil
	}

	var wrong []string
	for attempt := 0; attempt < verifyAttachmentAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(verifyBackoff(attempt))
		}
		attached, truncated, err := client.ReadAttachedPolicyPacks(target)
		if truncated {
			// Not retryable: a second read returns the same truncated page. For an attach, absence
			// proves nothing because the pack may sit on a page that was not returned; for a
			// detach, truncation can only hide a pack that is still attached.
			return nil
		}
		if IsTargetNotFound(err) {
			// The target is gone, so there is nothing to verify and no point retrying. Not an error:
			// the mutation was accepted, and a target that no longer exists takes its attachments
			// with it - which Exists reports on the next refresh.
			return nil
		}
		if err != nil {
			// Retryable: a blip on the confirmation read should cost one attempt, not abandon
			// verification entirely. If the target stays unreadable for the whole budget the guard
			// after the loop returns nil, preserving the rule that a mutation the server accepted is
			// not failed merely because the check could not run.
			continue
		}
		wrong = nil
		for _, pack := range packs {
			if policyPackInList(attached, pack) != wantAttached {
				wrong = append(wrong, pack)
			}
		}
		if len(wrong) == 0 {
			return nil
		}
	}

	if len(wrong) == 0 {
		// Every read failed, so nothing was ever compared - `wrong` is empty because the check never
		// ran, not because the write landed. Reporting a failure here would fail an apply that the
		// server accepted, which is exactly what the retry above exists to avoid.
		return nil
	}

	if wantAttached {
		return fmt.Errorf("error creating smart folder attachment: the API reported success but %v is not attached to %s. This usually means a concurrent attachment write to the same resource overwrote it; retrying the apply, or running with -parallelism=1, should converge", wrong, target)
	}
	return fmt.Errorf("error deleting smart folder attachment: the API reported success but %v is still attached to %s. This usually means a concurrent attachment write to the same resource overwrote the detach; retrying, or running with -parallelism=1, should converge", wrong, target)
}

// attachmentPacks normalises the smartFolders field, which the mutation accepts as either a
// single identifier or a list of them.
func attachmentPacks(input map[string]interface{}) []string {
	switch packs := input["smartFolders"].(type) {
	case string:
		return []string{packs}
	case []string:
		return packs
	case []interface{}:
		var out []string
		for _, pack := range packs {
			if s, ok := pack.(string); ok {
				out = append(out, s)
			}
		}
		return out
	}
	return nil
}
