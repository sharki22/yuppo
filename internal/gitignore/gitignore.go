package gitignore

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5/plumbing/format/gitignore"
)

func Load(rootPath string) (gitignore.Matcher, error) {
	gitignorePath := filepath.Join(rootPath, ".gitignore")

	_, err := os.Stat(gitignorePath)
	if os.IsNotExist(err) {
		return gitignore.NewMatcher(nil), nil
	}
	if err != nil {
		return nil, fmt.Errorf("stat .gitignore: %w", err)
	}

	bfs := osfs.New(rootPath)
	patterns, err := gitignore.ReadPatterns(bfs, nil)
	if err != nil {
		return nil, fmt.Errorf("read .gitignore: %w", err)
	}

	return gitignore.NewMatcher(patterns), nil
}
