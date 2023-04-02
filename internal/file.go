package internal

import (
	"fmt"
	"os"
	"strings"
)

// Writes a single checksum file containing all checksums in a single file
// Format: <algorithm>: <hash>
func WriteChecksumsContainerFile(path string, checksums []Checksum) error {
	var sb strings.Builder
	for _, c := range checksums {
		sb.WriteString(c.LineString() + "\n")
	}
	sb.WriteRune('\n')

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(sb.String())
	if err != nil {
		return err
	}

	return nil
}

// Writes a checksum file for each hash algorithm
func WriteChecksumsFiles(path string, checksums []Checksum) error {

	for _, c := range checksums {

		f, err := os.Create(path + "." + c.AlgorithmString())
		if err != nil {
			return err
		}
		defer f.Close()

		// Write (md5|sha256)sum compatible checksum
		_, err = f.WriteString(c.Hash + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Writes a checksum file for each hash algorithm
func WriteChecksumCompatibleFiles(path string, checksums []Checksum) error {

	for _, c := range checksums {

		f, err := os.Create(path + "." + c.AlgorithmString())
		if err != nil {
			return err
		}
		defer f.Close()

		// Write (md5|sha256)sum compatible checksum
		_, err = f.WriteString(fmt.Sprintf("%s *%s\n", c.Hash, path))
		if err != nil {
			return err
		}
	}

	return nil
}
