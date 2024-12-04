package main

import (
	"log"
	"os"
)

func main() {

	InitConfigStruct()
	InitLayerConfig()
	InitImageConfig()
	InitGlobalImageInfos()

	if len(os.Args) < 2 {
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar>\nother cmd is not supported for now!")
	}

	cmd := os.Args[1]

	switch cmd {
	case "create":
		if err := ecsCreate(os.Args[2]); err != nil {
			log.Println("Error generating RootFS:", err)
		}
	case "save":
		if err := ecsSave(os.Args[2], os.Args[3]); err != nil {
			log.Println("Error saving Image:", err)
		}
	default:
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar>\nother cmd is not supported for now!")
		os.Exit(1)
	}
}
