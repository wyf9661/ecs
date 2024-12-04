package main

import (
	"encoding/json"
	"fmt"
	"log"
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

	log.Println("RootFS generated at", rootPath)
	return nil
}

func ecsCreateConfigJson(jsonPath string) error {

	jsonData, err := json.MarshalIndent(GlobalConfigStruct, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return err
	}

	// To do (optional) use "github.com/tidwall/sjson" package to customlize config.json when creating bundle
	// sjson.Set() , sjson.SetBytes() ... etc.

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		log.Println("Error writing file:", err)
		return err
	}

	log.Println("config.json created successfully.")

	return nil
}

func ecsCreate(rootPath string) error {

	_, err := os.Stat(rootPath)
	if !os.IsNotExist(err) {
		log.Printf("The bundle %s already exists and will be removed...\n", rootPath)
		err = os.RemoveAll(rootPath)
		if err != nil {
			log.Fatalf("The bundle removed faild, please retry.\n")
		}
		log.Println("The bundle has been removed.")
	}

	// create rootfs
	if err := ecsCreateRootfs(rootPath + "/rootfs"); err != nil {
		log.Println("ecsCreateRootfs failed")
	}
	// create config.json
	if err := ecsCreateConfigJson(rootPath + "/config.json"); err != nil {
		log.Println("ecsCreateConfigJson failed")
	}
	return nil
}
