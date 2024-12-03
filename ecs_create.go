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
func ecsCreateRootfs(rootPath string) error {
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

func ecsCreateConfigJson(rootPath string) error {
	var jsonData []byte = []byte(jsonTemplate)

	// To do (optional) use "github.com/tidwall/sjson" package to customlize config.json when creating bundle
	// sjson.Set() , sjson.SetBytes() ... etc.

	if err := os.WriteFile(rootPath+"/config.json", jsonData, 0644); err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	fmt.Println("config.json created successfully.")

	return nil
}

func ecsCreate(rootPath string) error {
	// create rootfs
	if err := ecsCreateRootfs(rootPath); err != nil {
		fmt.Println("ecsCreateRootfs failed")
	}
	// create config.json
	if err := ecsCreateConfigJson(rootPath); err != nil {
		fmt.Println("ecsCreateConfigJson failed")
	}
	return nil
}
