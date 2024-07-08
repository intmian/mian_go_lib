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

type GitCommit struct {
	Hash    string
	Author  string
	Date    string
	Message string
}

// CompareGitVersion 比较git库的版本并返回差异
func CompareGitVersion(repoPath, version string) (bool, []GitCommit, error) {
	currentVersion, err := GetGitVersion(repoPath)
	if err != nil {
		return false, nil, err
	}

	if currentVersion == version {
		return true, nil, nil
	}

	// 获取所有版本的差异和具体改动
	cmd := exec.Command("git", "log", "--pretty=format:%H###%ar###%cn###%s", version+".."+currentVersion)
	cmd.Dir = repoPath
	out, err := cmd.Output()
	if err != nil {
		return false, nil, err
	}

	versions := strings.Split(string(out), "\n")
	commits := make([]GitCommit, 0, len(versions))
	for _, v := range versions {
		if v == "" {
			continue
		}
		s := strings.Split(v, "###")
		if s == nil || len(s) < 4 {
			continue
		}
		commits = append(commits, GitCommit{
			Hash:    s[0],
			Date:    s[1],
			Author:  s[2],
			Message: s[3],
		})
	}
	return false, commits, nil
}
