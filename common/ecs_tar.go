package common

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func TarDirectory(tarFileName string, sourceDir string) error {

	tarFile, err := os.Create(tarFileName)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}
