package gitops

import (
	"fmt"

	"github.com/go-git/go-git/v5"
)

func AutoCommit(repoPath, message string, autoPush bool) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("open repo: %w", err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("get worktree: %w", err)
	}

	if _, err := w.Add("."); err != nil {
		return fmt.Errorf("git add: %w", err)
	}

	if _, err := w.Commit(message, &git.CommitOptions{}); err != nil {
		return fmt.Errorf("git commit: %w", err)
	}

	if autoPush {
		if err := Push(repoPath); err != nil {
			return fmt.Errorf("git push: %w", err)
		}
	}

	return nil
}

func Push(repoPath string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("open repo for push: %w", err)
	}

	if err := repo.Push(&git.PushOptions{RemoteName: "origin"}); err != nil {
		return fmt.Errorf("push to origin: %w", err)
	}

	return nil
}
