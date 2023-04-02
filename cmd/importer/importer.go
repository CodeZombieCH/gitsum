package main

import (
	"bufio"
	"crypto"
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/codezombiech/gitsum/internal"
)

func main() {
	var hashFilePath string
	var repoPath string

	flag.StringVar(&hashFilePath, "hashfile", "", "path to hash file")
	flag.StringVar(&repoPath, "repo", "", "path to repo directory")

	flag.Parse()

	checkDirectoryPath("repo", repoPath)

	// read file line by line
	// normalize to unix path
	// create corresponding file in repo

	file, err := os.Open(hashFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// parse line
		if len(line) == 0 {
			continue
		}
		if line[:2] == "//" {
			continue
		}
		md5Hash := line[:32]
		localPath := line[33:]

		localPath = strings.ReplaceAll(localPath, "\\", "/")
		log.Printf("%s %s", md5Hash, localPath)

		if err := writeChecksums(path.Join(repoPath, localPath), md5Hash); err != nil {
			log.Fatalf("%v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func writeChecksums(filePath string, md5Hash string) error {
	if err := os.MkdirAll(path.Dir(filePath), 0775); err != nil {
		return err
	}

	checksums := []internal.Checksum{{Algorithm: crypto.MD5, Hash: md5Hash}}
	if err := internal.WriteChecksumsFiles(filePath, checksums); err != nil {
		return err
	}

	return nil
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
