package misc

import (
	"github.com/go-git/go-git/v5"
	"os/exec"
	"strings"
)

// GetGitVersion 获取git库的当前版本
func GetGitVersion(repoPath string) (string, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

// CompareGitVersion 比较git库的版本并返回差异
func CompareGitVersion(repoPath, version string) (bool, []string, error) {
	currentVersion, err := GetGitVersion(repoPath)
	if err != nil {
		return false, nil, err
	}

	if currentVersion == version {
		return true, nil, nil
	}

	// 获取所有版本的差异
	cmd := exec.Command("git", "log", "--pretty=format:%H", version+".."+currentVersion)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return false, nil, err
	}

	versions := strings.Split(string(out), "\n")
	return false, versions, nil
}
