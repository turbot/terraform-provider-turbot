package apiClient

import (
	"fmt"
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
var attachmentTargetLocks sync.Map // target identifier -> *sync.Mutex

// lockAttachmentTarget serialises attachment writes for a single target. It returns the unlock
// func so callers can defer it.
func lockAttachmentTarget(target string) func() {
	value, _ := attachmentTargetLocks.LoadOrStore(target, &sync.Mutex{})
	mutex := value.(*sync.Mutex)
	mutex.Lock()
	return mutex.Unlock
}

// attachmentTarget pulls the target identifier out of a mutation input, which is the value the
// lock is keyed on. Returns "" when absent, in which case the caller skips locking rather than
// serialising every attachment in the process behind one key.
func attachmentTarget(input map[string]interface{}) string {
	target, ok := input["resource"].(string)
	if !ok {
		return ""
	}
	return target
}

// verifyAttachmentAttempts / verifyAttachmentDelay bound the post-write confirmation below. The
// attachment list is not immediately consistent after a write, so a single read can report a
// successful attachment as missing; these allow a short settle before concluding it was lost.
const (
	verifyAttachmentAttempts = 4
	verifyAttachmentDelay    = 2 * time.Second
)

func (client *Client) CreateSmartFolderAttachment(input map[string]interface{}) (*TurbotResourceMetadata, error) {
	query := createSmartFolderAttachmentMutation()
	responseData := &CreateSmartFolderAttachResponse{}

	variables := map[string]interface{}{
		"input": input,
	}

	if target := attachmentTarget(input); target != "" {
		defer lockAttachmentTarget(target)()
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "smart folder attachment")
	}

	// The mutation can report success without the attachment persisting (see the note on
	// attachmentTargetLocks). Locking prevents that within one process, but a target addressed by
	// id in one place and by aka in another takes two different locks, and nothing serialises
	// separate terraform runs. Confirm the attachment landed so a lost write fails loudly here
	// rather than silently reappearing as drift on a later plan.
	if err := client.verifyAttachment(input); err != nil {
		return nil, err
	}

	return &responseData.SmartFolderAttach.Turbot, nil
}

// verifyAttachment confirms every policy pack named in input is attached to the target, retrying
// briefly to allow for read-after-write lag. A nil return means confirmed; verification that
// cannot run at all (unreadable input, or a target the caller cannot read) is not treated as
// failure, since the mutation itself already succeeded.
func (client *Client) verifyAttachment(input map[string]interface{}) error {
	target := attachmentTarget(input)
	if target == "" {
		return nil
	}
	wanted := attachmentPacks(input)
	if len(wanted) == 0 {
		return nil
	}

	var missing []string
	for attempt := 0; attempt < verifyAttachmentAttempts; attempt++ {
		if attempt > 0 {
			time.Sleep(verifyAttachmentDelay)
		}
		attached, truncated, err := client.ReadAttachedPolicyPacks(target)
		if err != nil || truncated {
			// Cannot verify - a read failure, or a truncated list in which absence proves
			// nothing. Do not fail a mutation that the server accepted.
			return nil
		}
		missing = nil
		for _, pack := range wanted {
			if !policyPackInList(attached, pack) {
				missing = append(missing, pack)
			}
		}
		if len(missing) == 0 {
			return nil
		}
	}

	return fmt.Errorf("error creating smart folder attachment: the API reported success but %v is not attached to %s. This usually means a concurrent attachment to the same resource overwrote it; retrying the apply, or running with -parallelism=1, should converge", missing, target)
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

func (client *Client) DeleteSmartFolderAttachment(input map[string]interface{}) error {
	query := detachSmartFolderAttachment()
	var responseData interface{}

	variables := map[string]interface{}{
		"input": input,
	}

	// Detach is the same read-modify-write on the target's list, so it takes the same lock.
	if target := attachmentTarget(input); target != "" {
		defer lockAttachmentTarget(target)()
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting smart folder attachment: %s", err.Error())
	}
	return nil
}
