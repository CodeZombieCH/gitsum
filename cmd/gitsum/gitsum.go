package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"flag"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/codezombiech/gitsum/internal"
)

func main() {

	var storePath string
	var repoPath string

	flag.StringVar(&storePath, "store", "", "path to store directory")
	flag.StringVar(&repoPath, "repo", "", "path to repo directory")

	flag.Parse()

	// Check args
	checkDirectoryPath("store", storePath)
	checkDirectoryPath("repo", repoPath)

	log.Printf("store: %s, repo: %s", storePath, repoPath)

	if err := sync(storePath, repoPath); err != nil {
		log.Fatalf("failed: %v", err)
	}
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

func sync(storePath string, repoPath string) error {
	// Delete repo workspace
	if err := resetChecksums(repoPath); err != nil {
		return err
	}

	// Create all checksum files
	if err := calcChecksums(storePath, repoPath); err != nil {
		return err
	}

	return nil
}

func resetChecksums(repoPath string) error {
	repoFS := os.DirFS(repoPath)

	err := fs.WalkDir(repoFS, ".", func(localPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("resetChecksums: failed to walk directory: %w", err)
		}
		if d.IsDir() {
			if localPath == "." {
				// Do not delete the repo directory itself
				return nil
			}
			if localPath == ".git" {
				// Do not delete the git repository
				return fs.SkipDir
			}
			if localPath == ".gitattributes" {
				// Do not delete the .gitattributes file
				return fs.SkipDir
			}

			err := os.RemoveAll(path.Join(repoPath, localPath))
			if err != nil {
				return fmt.Errorf("resetChecksums: failed to remove directory: %w", err)
			}
			return fs.SkipDir
		}

		if err := os.Remove(path.Join(repoPath, localPath)); err != nil {
			return fmt.Errorf("resetChecksums: failed to remove file: %w", err)
		}

		return nil
	})

	return err
}

func calcChecksums(storePath string, repoPath string) error {
	fileSystem := os.DirFS(storePath)

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
				// We do not wan't a `.git` directory in our checksum repo and overwrite it when generating checksums
				log.Printf("warning: ignoring git repository %v", path.Join(storePath, localPath))
				return fs.SkipDir
			}

			// Prepare directory in repo
			os.MkdirAll(path.Join(repoPath, localPath), 0775)
			return nil
		}

		checksums, err := checksumFile(path.Join(storePath, localPath))
		if err != nil {
			return fmt.Errorf("calcChecksums: failed to calculate checksums: %w", err)
		}

		if err := internal.WriteChecksumsFiles(path.Join(repoPath, localPath), checksums); err != nil {
			log.Printf("failed to write checksum file: %v", err)
		}
		return nil
	})

	return err
}

func checksumFile(path string) ([]internal.Checksum, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	algorithms := map[crypto.Hash]hash.Hash{
		crypto.MD5:    md5.New(),
		crypto.SHA256: sha256.New(),
		crypto.SHA512: sha512.New(),
	}

	checksums := make([]internal.Checksum, 0)
	for a, h := range algorithms {
		f.Seek(0, io.SeekStart)
		if _, err := io.Copy(h, f); err != nil {
			return nil, err
		}
		checksum := internal.Checksum{Algorithm: a, Hash: fmt.Sprintf("%x", h.Sum(nil))}
		checksums = append(checksums, checksum)
		log.Printf("%s: %s %s\n", path, checksum.Algorithm, checksum.Hash)
	}

	return checksums, nil
}
