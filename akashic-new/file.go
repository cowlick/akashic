package main

import (
	"io"
	"os"
	"path/filepath"
)

func copyFile(from string, to string) error {
	original, err := os.Open(from)
	if err != nil {
		return err
	}
	defer original.Close()

	target, err := os.Create(to)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, original)
	if err != nil {
		return err
	}

	return target.Sync()
}

func copyFiles(targetDir string, outDir string) error {
	return filepath.Walk(targetDir,
		func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			rel, err := filepath.Rel(targetDir, path)
			if err != nil {
				return err
			}
			outPath := filepath.Join(outDir, rel)

			if info.IsDir() {
				if outPath == outDir {
					return nil
				}
				return os.Mkdir(outPath, info.Mode())
			}

			return copyFile(path, outPath)
		})
}
