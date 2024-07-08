package misc

import "testing"

func TestGit(t *testing.T) {
	version, err := GetGitVersion("..\\..\\..")
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Fatal("GetGitVersion failed")
	}
	version = "fde32a8ce9c84f45ea36ba42b21b880f19257360"
	same, versions, err := CompareGitVersion("..\\..\\..", version)
	if err != nil {
		t.Fatal(err)
	}
	if !same || len(versions) != 0 {
		t.Fatalf("CompareGitVersion failed, same: %v, versions: %v", same, versions)
	}
}
