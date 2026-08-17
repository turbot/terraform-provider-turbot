package apiClient

import (
	"sync"
	"testing"

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
