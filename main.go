package main

import (
	"ecs/cmd"
	"ecs/common"
	"ecs/parser"
	"log"
	"os"
	"strings"
)

func main() {

	parser.InitConfigStruct()
	parser.InitLayerConfig()
	parser.InitImageConfig()
	parser.InitGlobalImageInfos()

	if len(os.Args) <= 2 {
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar:version>/<image.tar>\n\tecs genbase <basePath>")
	}

	input := os.Args[1]

	switch input {
	case "create":
		if err := cmd.EcsCreate(os.Args[2]); err != nil {
			log.Println("Error generating RootFS:", err)
		}
	case "genbase":
		if err := common.UntarGz("/opt/realevo-linux-tools/base/default.tar.gz", os.Args[2]); err != nil {
			log.Println("Error generating Base:", err)
		}
	case "save":

		if len(os.Args) <= 3 {
			log.Fatalf("Usage:\n\tecs save <bundle> <image.tar:version>/<image.tar>")
		}

		parts := strings.Split(os.Args[3], ":")
		var name, version string

		if len(parts) > 1 {
			name = strings.TrimSpace(parts[0])
			version = strings.TrimSpace(parts[1])
			log.Printf("set container tarball version with %s.", version)
		} else {
			name = os.Args[3]
			version = "latest"
			log.Println("No version configuration, `latest` will be set as default.")
		}

		if err := cmd.EcsSave(os.Args[2], name, version); err != nil {
			log.Println("Error saving Image:", err)
		}
	default:
		log.Fatalf("Usage:\n\tecs create <bundle>\n\tecs save <bundle> <image.tar:version>/<image.tar>\n\tecs genbase <basePath>")
	}
}
