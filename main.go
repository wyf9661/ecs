package main

import (
	"fmt"
	"os"
)

// create dir. If exists, ignored.
func createDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// create file. If exists, cover it.
func createFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	return err
}

// generate rootfs with sylixos directioy structure
func ecsCreate(rootPath string) error {
	// subdir lists
	subDirs := []string{
		"apps",
		"bin",
		"boot",
		"dev",
		"etc",
		"home",
		"lib",
		"proc",
		"qt",
		"root",
		"sbin",
		"tmp",
		"usr",
		"var",
	}

	// create subdir
	for _, dir := range subDirs {
		if err := createDir(rootPath + "/" + dir); err != nil {
			return err
		}
	}

	// create startup.sh and set stack size with 200000 as default
	if err := createFile(rootPath+"/etc/startup.sh", "shstack 200000\n"); err != nil {
		return err
	}

	fmt.Println("RootFS generated at", rootPath)
	return nil
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ecs create <rootPath>")
		os.Exit(1)
	}

	cmd := os.Args[1]
	// only support create now
	if cmd != "create" {
		fmt.Println("Usage: ecs create <rootPath> \n	other cmd is not supported for now!")
		os.Exit(1)
	}

	rootPath := os.Args[2] + "/rootfs"
	if err := ecsCreate(rootPath); err != nil {
		fmt.Println("Error generating RootFS:", err)
	}
}
