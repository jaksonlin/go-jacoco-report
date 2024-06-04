package util

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func GetFilesByPostfix(path string, postfix string) ([]string, error) {
	result := make([]string, 0)
	// walk the path to get all files
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), postfix) {
			result = append(result, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
