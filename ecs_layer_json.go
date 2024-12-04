package main

type ContainerConfig struct {
	Hostname     string                 `json:"Hostname"`
	Domainname   string                 `json:"Domainname"`
	User         string                 `json:"User"`
	AttachStdin  bool                   `json:"AttachStdin"`
	AttachStdout bool                   `json:"AttachStdout"`
	AttachStderr bool                   `json:"AttachStderr"`
	Tty          bool                   `json:"Tty"`
	OpenStdin    bool                   `json:"OpenStdin"`
	StdinOnce    bool                   `json:"StdinOnce"`
	Env          []string               `json:"Env"`
	Cmd          []string               `json:"Cmd"`
	Image        string                 `json:"Image"`
	Volumes      map[string]interface{} `json:"Volumes"`
	WorkingDir   string                 `json:"WorkingDir"`
	Entrypoint   []string               `json:"Entrypoint"`
	OnBuild      []string               `json:"OnBuild"`
	Labels       map[string]string      `json:"Labels"`
}

type LayerConfig struct {
	Id              string          `json:"id"`
	Created         string          `json:"created"`
	ContainerConfig ContainerConfig `json:"container_config"`
	Os              string          `json:"os"`
}

var GlobalLayerConfig *LayerConfig

func InitLayerConfig() {
	GlobalLayerConfig = &LayerConfig{
		Id:      "",
		Created: "",
		ContainerConfig: ContainerConfig{
			Hostname:     "",
			Domainname:   "",
			User:         "",
			AttachStdin:  false,
			AttachStdout: false,
			AttachStderr: false,
			Tty:          false,
			OpenStdin:    false,
			StdinOnce:    false,
			Env:          []string{},
			Cmd:          []string{},
			Image:        "",
			Volumes:      nil,
			WorkingDir:   "",
			Entrypoint:   []string{},
			OnBuild:      []string{},
			Labels:       nil,
		},
		Os: "sylixos",
	}
}
