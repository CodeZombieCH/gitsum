package internal

import (
	"crypto"
	"fmt"
	"strconv"
	"strings"
)

type Checksum struct {
	Algorithm crypto.Hash
	Hash      string
}

func (c Checksum) AlgorithmString() string {
	return AlgorithmString(c.Algorithm)
}

func (c *Checksum) LineString() string {
	return fmt.Sprintf("%s: %s\n", c.AlgorithmString(), c.Hash)
}

func AlgorithmString(h crypto.Hash) string {
	switch h {
	case crypto.MD5:
		return "md5"
	case crypto.SHA256:
		return "sha256"
	case crypto.SHA512:
		return "sha512"
	default:
		return "unknown hash value " + strconv.Itoa(int(h))

	}
}

func ParseAlgorithm(raw string) (crypto.Hash, error) {
	switch strings.ToLower(raw) {
	case "md5":
		return crypto.MD5, nil
	case "sha256":
		return crypto.SHA256, nil
	case "sha512":
		return crypto.SHA512, nil
	default:
		return 0, fmt.Errorf("unknown hash value " + raw)
	}
}
