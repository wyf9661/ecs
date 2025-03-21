package parser

type Platform struct {
	Os   string `json:"os"`
	Arch string `json:"arch"`
}

type Process struct {
	User User     `json:"user"`
	Args []string `json:"args"`
	Env  []string `json:"env"`
	Cwd  string   `json:"cwd"`
}

type User struct {
	Uid int `json:"uid"`
	Gid int `json:"gid"`
}

type Root struct {
	Path     string `json:"path"`
	Readonly bool   `json:"readonly"`
}

type Mount struct {
	Destination string   `json:"destination"`
	Source      string   `json:"source"`
	Options     []string `json:"options"`
}

type Sylixos struct {
	Devices   []Device  `json:"devices"`
	Resources Resources `json:"resources"`
	Commands  []string  `json:"commands"`
	Network   Network   `json:"network"`
}

type Device struct {
	Path   string `json:"path"`
	Access string `json:"access"`
}

type Resources struct {
	Cpu          Cpu           `json:"cpu"`
	Itimer       Itimer        `json:"itimer"`
	Affinity     []interface{} `json:"affinity"`
	Memory       Memory        `json:"memory"`
	KernelObject KernelObject  `json:"kernelObject"`
	Disk         Disk          `json:"disk"`
}

type Cpu struct {
	HighestPrio int `json:"highestPrio"`
	LowestPrio  int `json:"lowestPrio"`
	DefaultPrio int `json:"defaultPrio"`
}

type Itimer struct {
	DefaultPrio int `json:"defaultPrio"`
}

type Memory struct {
	KheapLimit    int `json:"kheapLimit"`
	MemoryLimitMB int `json:"memoryLimitMB"`
}

type KernelObject struct {
	ThreadLimit        int `json:"threadLimit"`
	ThreadPoolLimit    int `json:"threadPoolLimit"`
	EventLimit         int `json:"eventLimit"`
	EventSetLimit      int `json:"eventSetLimit"`
	PartitionLimit     int `json:"partitionLimit"`
	RegionLimit        int `json:"regionLimit"`
	MsgQueueLimit      int `json:"msgQueueLimit"`
	TimerLimit         int `json:"timerLimit"`
	RmsLimit           int `json:"rmsLimit"`
	ThreadVarLimit     int `json:"threadVarLimit"`
	PosixMqueueLimit   int `json:"posixMqueueLimit"`
	DlopenLibraryLimit int `json:"dlopenLibraryLimit"`
	XsiipcLimit        int `json:"xsiipcLimit"`
	SocketLimit        int `json:"socketLimit"`
	SrtpLimit          int `json:"srtpLimit"`
	DeviceLimit        int `json:"deviceLimit"`
}

type Disk struct {
	LimitMB int `json:"limitMB"`
}

type Network struct {
	TelnetdEnable bool `json:"telnetdEnable"`
	FtpdEnable    bool `json:"ftpdEnable"`
	SshdEnable    bool `json:"sshdEnable"`
}

type ConfigStruct struct {
	OciVersion string   `json:"ociVersion"`
	Platform   Platform `json:"platform"`
	Process    Process  `json:"process"`
	Root       Root     `json:"root"`
	Hostname   string   `json:"hostname"`
	Mounts     []Mount  `json:"mounts"`
	Sylixos    Sylixos  `json:"sylixos"`
}

var GlobalConfigStruct *ConfigStruct

func InitConfigStruct() {
	GlobalConfigStruct = &ConfigStruct{
		OciVersion: "1.0.0",
		Platform: Platform{
			Os:   "sylixos",
			Arch: "noarch",
		},
		Process: Process{
			User: User{
				Uid: 0,
				Gid: 0,
			},
			Args: []string{""},
			Env: []string{
				"PATH=/usr/bin:/bin:/usr/pkg/sbin:/sbin:/usr/local/bin",
				"LD_LIBRARY_PATH=/usr/lib:/lib:/usr/local/lib",
			},
			Cwd: "/",
		},
		Root: Root{
			Path:     "rootfs",
			Readonly: false,
		},
		Hostname: "sylixos_ecs",
		Mounts: []Mount{
			{
				Destination: "/lib",
				Source:      "/lib",
				Options:     []string{"ro"},
			},
			{
				Destination: "/bin/uname",
				Source:      "/bin/uname",
				Options:     []string{"rx"},
			},
			{
				Destination: "/bin/tar",
				Source:      "/bin/tar",
				Options:     []string{"ro"},
			},
		},
		Sylixos: Sylixos{
			Devices: []Device{
				{
					Path:   "/dev/fb0",
					Access: "rw",
				},
				{
					Path:   "/dev/rjgt102",
					Access: "rw",
				},
				{
					Path:   "/dev/lkt4304",
					Access: "rw",
				},
				{
					Path:   "/dev/input/xmse",
					Access: "rw",
				},
				{
					Path:   "/dev/input/xkbd",
					Access: "rw",
				},
				{
					Path:   "/dev/net/vnd",
					Access: "rw",
				},
			},
			Resources: Resources{
				Cpu: Cpu{
					HighestPrio: 150,
					LowestPrio:  250,
					DefaultPrio: 200,
				},
				Itimer: Itimer{
					DefaultPrio: 200,
				},
				Affinity: []interface{}{},
				Memory: Memory{
					KheapLimit:    536870912,
					MemoryLimitMB: 2048,
				},
				KernelObject: KernelObject{
					ThreadLimit:        4096,
					ThreadPoolLimit:    32,
					EventLimit:         32768,
					EventSetLimit:      500,
					PartitionLimit:     6000,
					RegionLimit:        50,
					MsgQueueLimit:      8192,
					TimerLimit:         64,
					RmsLimit:           32,
					ThreadVarLimit:     16,
					PosixMqueueLimit:   300,
					DlopenLibraryLimit: 50,
					XsiipcLimit:        100,
					SocketLimit:        1024,
					SrtpLimit:          30,
					DeviceLimit:        60,
				},
				Disk: Disk{
					LimitMB: 2048,
				},
			},
			Commands: []string{
				"exec",
				"top",
				"cpuus",
				"vi",
				"cat",
				"touch",
				"ps",
				"ts",
				"tp",
				"ss",
				"ints",
				"ls",
				"cd",
				"pwd",
				"modules",
				"varload",
				"varsave",
				"shstack",
				"srtp",
				"shfile",
				"help",
				"debug",
				"shell",
				"ll",
				"sync",
				"ln",
				"kill",
				"free",
				"ifconfig",
				"mems",
				"env",
				"rm",
				"exit",
				"mkdir",
				"echo",
				"mv",
			},
			Network: Network{
				TelnetdEnable: true,
				FtpdEnable:    true,
				SshdEnable:    false,
			},
		},
	}
}
