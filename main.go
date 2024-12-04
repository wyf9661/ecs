package main

import (
	"ecs/cmd"
	"ecs/parser"
	"log"
	"os"
)

func main() {

	parser.InitConfigStruct()
	parser.InitLayerConfig()
	parser.InitImageConfig()
	parser.InitGlobalImageInfos()

	if len(os.Args) < 2 {
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar>\nother cmd is not supported for now!")
	}

	input := os.Args[1]

	switch input {
	case "create":
		if err := cmd.EcsCreate(os.Args[2]); err != nil {
			log.Println("Error generating RootFS:", err)
		}
	case "save":
		if err := cmd.EcsSave(os.Args[2], os.Args[3]); err != nil {
			log.Println("Error saving Image:", err)
		}
	default:
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar>\nother cmd is not supported for now!")
		os.Exit(1)
	}
}
