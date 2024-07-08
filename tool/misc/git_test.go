package misc

import "testing"

func TestGit(t *testing.T) {
	version, err := GetGitVersion("..\\..")
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Fatal("GetGitVersion failed")
	}
	version = "0c5d1543d97312ebbfb3db54c5fdaaceec4b1d24"
	same, versions, err := CompareGitVersion("..\\..", version)

	if err != nil {
		t.Fatal(err)
	}
	if same || len(versions) == 0 {
		t.Fatalf("CompareGitVersion failed, same: %v, versions: %v", same, versions)
	}
}
