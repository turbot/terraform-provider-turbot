package apiClient

import (
	"sync"
	"testing"
	"time"

	"github.com/machinebox/graphql"
	"github.com/stretchr/testify/assert"
)

// The lock is keyed on the target, so attachments to the same resource serialise while
// attachments to different resources still run concurrently.
func TestAttachmentTarget(t *testing.T) {
	var tests = []struct {
		name     string
		input    map[string]interface{}
		expected string
	}{
		{"target present", map[string]interface{}{"resource": "191926035367605", "smartFolders": "123"}, "191926035367605"},
		{"target as aka", map[string]interface{}{"resource": "my_bu_folder"}, "my_bu_folder"},
		{"target absent", map[string]interface{}{"smartFolders": "123"}, ""},
		{"target wrong type", map[string]interface{}{"resource": 12345}, ""},
		{"empty input", map[string]interface{}{}, ""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, attachmentTarget(test.input))
		})
	}
}

// The mutation accepts smartFolders as a bare identifier or a list; verification has to read
// both shapes or it would silently check nothing.
func TestAttachmentPacks(t *testing.T) {
	var tests = []struct {
		name     string
		input    map[string]interface{}
		expected []string
	}{
		{"single string - what both resources send", map[string]interface{}{"smartFolders": "343645598782139"}, []string{"343645598782139"}},
		{"string slice", map[string]interface{}{"smartFolders": []string{"a", "b"}}, []string{"a", "b"}},
		{"interface slice", map[string]interface{}{"smartFolders": []interface{}{"a", "b"}}, []string{"a", "b"}},
		{"interface slice with non-strings", map[string]interface{}{"smartFolders": []interface{}{"a", 7}}, []string{"a"}},
		{"absent", map[string]interface{}{"resource": "x"}, nil},
		{"wrong type", map[string]interface{}{"smartFolders": 42}, nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, attachmentPacks(test.input))
		})
	}
}

// Two writers for the SAME target must not overlap - that overlap is what loses attachments
// server-side. Asserted by counting concurrent holders rather than by timing.
func TestLockAttachmentTargetSerialisesSameTarget(t *testing.T) {
	var (
		mu       sync.Mutex
		inside   int
		maxSeen  int
		waitAll  sync.WaitGroup
		iterates = 50
	)

	for i := 0; i < iterates; i++ {
		waitAll.Add(1)
		go func() {
			defer waitAll.Done()
			unlock := lockAttachmentTarget("same-target")
			defer unlock()

			mu.Lock()
			inside++
			if inside > maxSeen {
				maxSeen = inside
			}
			mu.Unlock()

			mu.Lock()
			inside--
			mu.Unlock()
		}()
	}
	waitAll.Wait()

	assert.Equal(t, 1, maxSeen, "writes to the same target must never overlap")
}

// Different targets must not block each other, otherwise a large apply serialises entirely.
func TestLockAttachmentTargetAllowsDifferentTargets(t *testing.T) {
	first := lockAttachmentTarget("target-a")
	defer first()

	acquired := make(chan struct{})
	go func() {
		unlock := lockAttachmentTarget("target-b")
		defer unlock()
		close(acquired)
	}()

	// target-b must be obtainable while target-a is held; a deadlock here means the lock is
	// global rather than per target.
	<-acquired
}

// Re-locking the same target reuses the stored mutex rather than creating a fresh one, which is
// what makes the serialisation real.
func TestLockAttachmentTargetReusesMutexPerTarget(t *testing.T) {
	unlock := lockAttachmentTarget("reuse-target")
	firstValue, ok := attachmentTargetLocks.Load("reuse-target")
	assert.True(t, ok, "lock must be stored for the target")
	unlock()

	unlock = lockAttachmentTarget("reuse-target")
	secondValue, _ := attachmentTargetLocks.Load("reuse-target")
	unlock()

	assert.True(t, firstValue == secondValue, "the same target must map to the same mutex")
}

// The backoff ramps rather than sitting flat, because the wait is held under the target's lock:
// a coarse first delay would queue every same-target write behind it.
func TestVerifyBackoffRamps(t *testing.T) {
	var tests = []struct {
		attempt  int
		expected time.Duration
	}{
		{1, 250 * time.Millisecond},
		{2, 500 * time.Millisecond},
		{3, 1000 * time.Millisecond},
		{4, 2000 * time.Millisecond},
	}

	var total time.Duration
	for _, test := range tests {
		assert.Equal(t, test.expected, verifyBackoff(test.attempt), "attempt %d", test.attempt)
		total += test.expected
	}

	// verifyAttachmentAttempts attempts means verifyAttachmentAttempts-1 sleeps; keep the whole
	// settle window bounded so a stuck verification cannot hold a target's lock indefinitely.
	assert.Equal(t, 4, verifyAttachmentAttempts-1, "backoff table must cover every retry")
	assert.Equal(t, 3750*time.Millisecond, total, "total settle window")
}

// The shift in verifyBackoff is negative for attempt 0, which panics. The current call site guards
// against it, so this is a latent-panic guard rather than a live bug - but a panic in a provider is
// one refactor away.
func TestVerifyBackoffDoesNotPanicBelowFirstRetry(t *testing.T) {
	for _, attempt := range []int{0, -1, -100} {
		assert.NotPanics(t, func() { verifyBackoff(attempt) }, "attempt %d", attempt)
		assert.Equal(t, time.Duration(0), verifyBackoff(attempt), "attempt %d", attempt)
	}
}

// The resolved lock key must be cached, so attaching N packs to one target resolves it once rather
// than N times, and so every writer agrees on the key even if one resolution fails.
func TestAttachmentLockKeyIsCachedPerTarget(t *testing.T) {
	attachmentLockKeys.Store("cached-target", "999888777")
	defer attachmentLockKeys.Delete("cached-target")

	// A nil client would panic if the cache were bypassed, which is what makes this a real
	// assertion that the cached value short-circuits the API read.
	var client *Client
	assert.NotPanics(t, func() {
		assert.Equal(t, "999888777", client.attachmentLockKey("cached-target"))
	}, "a cached key must not trigger a resource read")
}

// Two writers naming the same target differently - one by id, one by aka - must end up on the same
// lock. That is the whole reason the key is resolved rather than used raw.
func TestAttachmentLockKeyUnifiesIdAndAka(t *testing.T) {
	attachmentLockKeys.Store("391406345032847", "391406345032847")
	attachmentLockKeys.Store("my_folder_aka", "391406345032847")
	defer func() {
		attachmentLockKeys.Delete("391406345032847")
		attachmentLockKeys.Delete("my_folder_aka")
	}()

	var client *Client
	byId := client.attachmentLockKey("391406345032847")
	byAka := client.attachmentLockKey("my_folder_aka")
	assert.Equal(t, byId, byAka, "id and aka for one target must resolve to the same lock key")

	// and that shared key must therefore serialise them
	unlock := lockAttachmentTarget(byId)
	acquired := make(chan struct{})
	go func() {
		u := lockAttachmentTarget(byAka)
		defer u()
		close(acquired)
	}()
	select {
	case <-acquired:
		t.Fatal("id and aka writers acquired the lock concurrently")
	case <-time.After(50 * time.Millisecond):
		// correct: blocked behind the first holder
	}
	unlock()
	<-acquired
}

// The read-failure branch had no coverage, which is how a spurious failure lived in it: with every
// read erroring, the loop exhausted its budget and then reported "the API reported success but []
// is not attached", naming no packs because none were ever compared. An unreadable target must
// return nil - a mutation the server accepted is not failed because the check could not run.
func TestVerifyReturnsNilWhenTargetUnreadable(t *testing.T) {
	// shrink the budget so the test does not burn the real backoff twice
	defer func(attempts int, delay time.Duration) {
		verifyAttachmentAttempts, verifyAttachmentBaseDelay = attempts, delay
	}(verifyAttachmentAttempts, verifyAttachmentBaseDelay)
	verifyAttachmentAttempts, verifyAttachmentBaseDelay = 3, time.Millisecond

	// 127.0.0.1:1 refuses immediately, so every confirmation read errors
	client := &Client{Graphql: graphql.NewClient("http://127.0.0.1:1/graphql")}
	input := map[string]interface{}{"resource": "12345", "smartFolders": "678"}

	assert.NoError(t, client.verifyAttachmentState(input, true),
		"an unreadable target must not fail an attach the server accepted")
	assert.NoError(t, client.verifyAttachmentState(input, false),
		"an unreadable target must not fail a detach the server accepted")
}
