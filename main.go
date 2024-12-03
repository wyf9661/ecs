package main

import (
	"fmt"
	"os"
)

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

	_, err := os.Stat(rootPath)
	if !os.IsNotExist(err) {
		fmt.Printf("The bundle %s already exists and will be removed...\n", rootPath)
		err = os.RemoveAll(rootPath)
		if err != nil {
			fmt.Printf("The bundle removed faild, please retry.\n", rootPath)
			os.Exit(1)
		}
		fmt.Println("The bundle has been removed.")
	}

	if err := ecsCreate(rootPath); err != nil {
		fmt.Println("Error generating RootFS:", err)
	}
}
