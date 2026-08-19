package history

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCleanupOldFiles checks the age cutoff and the owner-only permissions,
// the two things the hardened history package is responsible for.
func TestCleanupOldFiles(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	// A file inside the retention window and one well outside it.
	fresh := time.Now().Format("2006-01-02_15-04-05")
	stale := time.Now().Add(-2 * MaxHistoryAge).Format("2006-01-02_15-04-05")
	for _, ts := range []string{fresh, stale} {
		dir, err := EnsureHistoryDir()
		if err != nil {
			t.Fatalf("EnsureHistoryDir: %v", err)
		}
		name := ts + "_proj_plan.txt"
		if err := os.WriteFile(filepath.Join(dir, name), []byte("x"), filePerm); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	dir, _ := GetHistoryDir()
	if info, err := os.Stat(dir); err != nil {
		t.Fatalf("stat dir: %v", err)
	} else if perm := info.Mode().Perm(); perm != dirPerm {
		t.Errorf("history dir perm = %o, want %o", perm, dirPerm)
	}

	deleted, err := CleanupOldFiles()
	if err != nil {
		t.Fatalf("CleanupOldFiles: %v", err)
	}
	if deleted != 1 {
		t.Errorf("deleted = %d, want 1", deleted)
	}

	left, err := ListEntries("")
	if err != nil {
		t.Fatalf("ListEntries: %v", err)
	}
	if len(left) != 1 || !left[0].Timestamp.After(time.Now().Add(-MaxHistoryAge)) {
		t.Errorf("expected only the fresh entry to survive, got %d entries", len(left))
	}
	if info, err := os.Stat(left[0].Path); err != nil {
		t.Fatalf("stat file: %v", err)
	} else if perm := info.Mode().Perm(); perm != filePerm {
		t.Errorf("history file perm = %o, want %o", perm, filePerm)
	}
}
