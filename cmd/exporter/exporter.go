package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
)

func main() {
	var repoPath string
	var hashFilePath string
	var algorithm string

	flag.StringVar(&repoPath, "repo", "", "path to repo directory")
	flag.StringVar(&hashFilePath, "hashfile", "", "path to hash file")
	flag.StringVar(&algorithm, "algorithm", "", "hash algorithm")

	flag.Parse()

	checkDirectoryPath("repo", repoPath)
	if len(algorithm) <= 0 {
		log.Fatalf("invalid parameter \"algorithm\": %v", algorithm)
	}

	f, err := os.Create(hashFilePath)
	if err != nil {
		log.Fatalf("failed to create hash file: %v", err)
	}
	defer f.Close()

	writeChecksums(f, repoPath, algorithm)
}

func writeChecksums(f *os.File, repoPath string, algorithm string) error {
	fileSuffix := "." + algorithm
	fileSystem := os.DirFS(repoPath)

	err := fs.WalkDir(fileSystem, ".", func(localPath string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("warning: %v", err)
			return fs.SkipDir
		}

		if d.IsDir() {
			if localPath == "." {
				// Nothing to do with the store directory itself
				return nil
			}
			if localPath == ".git" {
				// Skip .git repository
				return fs.SkipDir
			}
			return nil
		}

		if path.Ext(localPath) != fileSuffix {
			return nil
		}

		// Read checksum file
		b, err := os.ReadFile(path.Join(repoPath, localPath))
		if err != nil {
			return fmt.Errorf("writeChecksums: failed to read checksum file: %w", err)
		}

		content := string(b)
		hash := content[:len(content)-1]

		// Write to checksums file
		actualFilePath := localPath[:len(localPath)-len(fileSuffix)]
		_, err = f.WriteString(fmt.Sprintf("%s *%s\n", hash, actualFilePath))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func checkDirectoryPath(argName string, path string) {
	if len(path) == 0 {
		log.Panicf("%s: empty path", argName)
	}

	f, err := os.Stat(path)
	if err != nil {
		log.Panicf("%s: invalid path or not accessible", argName)
	}

	if !f.IsDir() {
		log.Panicf("%s: not a directory", argName)
	}
}
