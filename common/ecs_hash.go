package common

import (
	"crypto/sha256"
	"io"
	"os"
)

func CalculateFileSHA256(filePath string) ([]byte, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
