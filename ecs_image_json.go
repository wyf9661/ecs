package main

type Rootfs struct {
	Type    string   `json:"type"`
	DiffIds []string `json:"diff_ids"`
}

type History struct {
	Created    string `json:"created"`
	CreatedBy  string `json:"created_by"`
	EmptyLayer bool   `json:"empty_layer"`
}

type Config struct {
	User       string            `json:"User"`
	Env        []string          `json:"Env"`
	Entrypoint []string          `json:"Entrypoint"`
	WorkingDir string            `json:"WorkingDir"`
	Labels     map[string]string `json:"Labels"`
}

type ImageConfig struct {
	Created      string    `json:"created"`
	Architecture string    `json:"architecture"`
	Os           string    `json:"os"`
	Rootfs       Rootfs    `json:"rootfs"`
	History      []History `json:"history"`
	Config       Config    `json:"config"`
}

var GlobalImageConfig *ImageConfig

func InitImageConfig() {
	GlobalImageConfig = &ImageConfig{
		Created:      "",
		Architecture: "",
		Os:           "sylixos",
		Rootfs: Rootfs{
			Type:    "layers",
			DiffIds: []string{"sha256:41162f3525f66c0c775e324038601280f0b3f2891add95d652707c2fe652b1fc"},
		},
		History: []History{
			{
				Created:    "",
				CreatedBy:  "ecs (go version for built)",
				EmptyLayer: true,
			},
		},
		Config: Config{
			User:       "1000:1000",
			Env:        []string{"PATH=/usr/bin:/bin:/usr/pkg/sbin:/sbin:/usr/local/bin"},
			Entrypoint: []string{"/apps/ecs_demo"},
			WorkingDir: "/rootfs",
			Labels: map[string]string{
				"hostname":                                          "sylixos_ecs",
				"io.buildah.version":                                "1.21.3",
				"mounts.0.destination":                              "/qt",
				"mounts.0.options":                                  "rx",
				"mounts.0.source":                                   "/qt",
				"sylixos.commands.0":                                "exec",
				"sylixos.commands.1":                                "top",
				"sylixos.devices.0.access":                          "rw",
				"sylixos.devices.0.path":                            "/dev/fb0",
				"sylixos.network.ftpdEnable":                        "true",
				"sylixos.network.telnetdEnable":                     "true",
				"sylixos.resources.cpu.highestPrio":                 "160",
				"sylixos.resources.cpu.lowestPrio":                  "250",
				"sylixos.resources.disk.limitMB":                    "2048",
				"sylixos.resources.kernelObject.deviceLimit":        "60",
				"sylixos.resources.kernelObject.dlopenLibraryLimit": "50",
				"sylixos.resources.kernelObject.eventLimit":         "800",
				"sylixos.resources.kernelObject.eventSetLimit":      "50",
				"sylixos.resources.kernelObject.msgQueueLimit":      "50",
				"sylixos.resources.kernelObject.partitionLimit":     "5",
				"sylixos.resources.kernelObject.posixMqueueLimit":   "300",
				"sylixos.resources.kernelObject.regionLimit":        "5",
				"sylixos.resources.kernelObject.rmsLimit":           "5",
				"sylixos.resources.kernelObject.socketLimit":        "50",
				"sylixos.resources.kernelObject.srtpLimit":          "10",
				"sylixos.resources.kernelObject.threadLimit":        "300",
				"sylixos.resources.kernelObject.threadPoolLimit":    "1",
				"sylixos.resources.kernelObject.threadVarLimit":     "2",
				"sylixos.resources.kernelObject.timerLimit":         "5",
				"sylixos.resources.kernelObject.xsiipcLimit":        "100",
				"sylixos.resources.memory.kheapLimit":               "2097152",
				"sylixos.resources.memory.memoryLimitMB":            "2048",
			},
		},
	}
}
