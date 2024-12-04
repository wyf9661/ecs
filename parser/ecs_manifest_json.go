package parser

type ImageInfo struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

var GlobalImageInfos []*ImageInfo

func InitGlobalImageInfos() {
	GlobalImageInfos = []*ImageInfo{
		{
			Config:   "",
			RepoTags: []string{""},
			Layers:   []string{""},
		},
	}
}
