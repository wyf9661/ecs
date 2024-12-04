package cmd

import (
	"ecs/common"
	"ecs/parser"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func ecsSaveCreateWorkDir(workDir string) error {
	log.Println("Creating working directory...")
	_, err := os.Stat(workDir)
	if !os.IsNotExist(err) {
		log.Println("working directory: ./ecs_temp already exists and will be removed...")
		err = os.RemoveAll(workDir)
		if err != nil {
			log.Fatalf("The bundle removed faild, please retry.\n")
			return err
		}
		log.Println("working directory has been removed.")
	}

	if err = common.CreateDir(workDir); err != nil {
		log.Println("Failed to create working directory")
		return err
	}

	return nil
}

func ecsSaveCreateLayerTar(bundlePath string, layImage string) error {
	if err := common.TarDirectory(layImage, bundlePath); err != nil {
		log.Println("Failed to create lay tar")
		return err
	}
	return nil
}

func ecsSaveCreateLayerJson(layerPath string, sha256Value string) error {

	parser.GlobalLayerConfig.Id = sha256Value

	now := time.Now()
	formattedTime := now.Format("2006-01-02T15:04:05-07:00")
	parser.GlobalLayerConfig.Created = formattedTime

	jsonData, err := json.MarshalIndent(parser.GlobalLayerConfig, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return err
	}

	if err = os.WriteFile(layerPath+"/json", []byte(jsonData), 0644); err != nil {
		log.Println("Error writing file:", err)
		return err
	}
	return nil
}

func ecsSaveCreateLayerVersion(layerPath string) error {
	if err := common.CreateFile(layerPath+"/VERSION", "1.0"); err != nil {
		return err
	}
	return nil
}

func ecsSaveCreateLayer(bundlePath string, workDir string) (string, error) {
	log.Println("Creating layer...")

	log.Println("Creating layer directory...")
	if err := common.CreateDir(workDir + "/latest"); err != nil {
		log.Println("Failed to create working directory")
		return "", err
	}

	log.Println("Creating layer tarball...")
	if err := ecsSaveCreateLayerTar(bundlePath, workDir+"/latest"+"/layer.tar"); err != nil {
		log.Println("Failed to create working directory")
		return "", err
	}

	sha256Value, err := common.CalculateFileSHA256(workDir + "/latest" + "/layer.tar")
	if err != nil {
		log.Println("Failed to calculate sha256")
		return "", err
	}

	sha256sums := fmt.Sprintf("%x", sha256Value)

	log.Printf("Creating layerJson, sha256 is %s", sha256sums)
	if err := ecsSaveCreateLayerJson(workDir+"/latest", sha256sums); err != nil {
		log.Println("Failed to create layer json")
		return "", err
	}

	log.Println("Creating layer version...")
	if err := ecsSaveCreateLayerVersion(workDir + "/latest"); err != nil {
		log.Println("Failed to create layer version")
		return "", err
	}

	log.Println("Rename layer directory...")
	if err := os.Rename(workDir+"/latest", workDir+"/"+sha256sums); err != nil {
		log.Printf("Error renaming folder: %v\n", err)
		return "", err
	}

	return sha256sums, nil
}

func ecsSaveImageConfig(bundlePath string, workDir string, layerSha256 string) (string, error) {

	now := time.Now()
	formattedTime := now.Format(time.RFC3339Nano)
	parser.GlobalImageConfig.Created = formattedTime

	parser.GlobalImageConfig.Architecture = parser.GlobalConfigStruct.Platform.Arch
	parser.GlobalImageConfig.Os = parser.GlobalConfigStruct.Platform.Os
	parser.GlobalImageConfig.Config.User = strconv.Itoa(parser.GlobalConfigStruct.Process.User.Uid) + ":" + strconv.Itoa(parser.GlobalConfigStruct.Process.User.Gid)
	parser.GlobalImageConfig.Config.Env = parser.GlobalConfigStruct.Process.Env
	parser.GlobalImageConfig.Config.Entrypoint = parser.GlobalConfigStruct.Process.Args
	parser.GlobalImageConfig.Config.WorkingDir = parser.GlobalConfigStruct.Process.Cwd

	parser.GlobalImageConfig.Config.Labels["hostname"] = parser.GlobalConfigStruct.Hostname

	for i := 0; i < len(parser.GlobalConfigStruct.Mounts); i++ {
		parser.GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".destination"] = parser.GlobalConfigStruct.Mounts[i].Destination
		parser.GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".source"] = parser.GlobalConfigStruct.Mounts[i].Source
		parser.GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".options"] = parser.GlobalConfigStruct.Mounts[i].Options[0]
	}

	for i := 0; i < len(parser.GlobalConfigStruct.Sylixos.Commands); i++ {
		parser.GlobalImageConfig.Config.Labels["sylixos.commands."+strconv.Itoa(i)] = parser.GlobalConfigStruct.Sylixos.Commands[i]
	}

	for i := 0; i < len(parser.GlobalConfigStruct.Sylixos.Devices); i++ {
		parser.GlobalImageConfig.Config.Labels["sylixos.devices."+strconv.Itoa(i)+".path"] = parser.GlobalConfigStruct.Sylixos.Devices[i].Path
		parser.GlobalImageConfig.Config.Labels["sylixos.devices."+strconv.Itoa(i)+".access"] = parser.GlobalConfigStruct.Sylixos.Devices[i].Access
	}

	parser.GlobalImageConfig.Config.Labels["sylixos.resources.cpu.highestPrio"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.Cpu.HighestPrio)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.cpu.lowestPrio"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.Cpu.LowestPrio)

	parser.GlobalImageConfig.Config.Labels["sylixos.resources.memory.kheapLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.Memory.KheapLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.memory.memoryLimitMB"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.Memory.MemoryLimitMB)

	parser.GlobalImageConfig.Config.Labels["sylixos.resources.disk.limitMB"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.Disk.LimitMB)

	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.deviceLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.DeviceLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.dlopenLibraryLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.DlopenLibraryLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.eventLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.EventLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.eventSetLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.EventSetLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.msgQueueLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.MsgQueueLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.partitionLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.PartitionLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.posixMqueueLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.PosixMqueueLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.regionLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.RegionLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.rmsLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.RmsLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.socketLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.SocketLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.srtpLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.SrtpLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadPoolLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadPoolLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadVarLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadVarLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.timerLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.TimerLimit)
	parser.GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.xsiipcLimit"] = strconv.Itoa(parser.GlobalConfigStruct.Sylixos.Resources.KernelObject.XsiipcLimit)

	parser.GlobalImageConfig.Config.Labels["sylixos.network.ftpdEnable"] = fmt.Sprintf("%t", parser.GlobalConfigStruct.Sylixos.Network.FtpdEnable)
	parser.GlobalImageConfig.Config.Labels["sylixos.network.telnetdEnable"] = fmt.Sprintf("%t", parser.GlobalConfigStruct.Sylixos.Network.TelnetdEnable)
	parser.GlobalImageConfig.Config.Labels["sylixos.network.sshdEnable"] = fmt.Sprintf("%t", parser.GlobalConfigStruct.Sylixos.Network.SshdEnable)

	parser.GlobalImageConfig.Rootfs.Type = "layers"
	parser.GlobalImageConfig.Rootfs.DiffIds[0] = "sha256:" + layerSha256

	parser.GlobalImageConfig.History[0].Created = formattedTime
	parser.GlobalImageConfig.History[0].CreatedBy = "ecs (go version for built)"
	parser.GlobalImageConfig.History[0].EmptyLayer = false

	jsonData, err := json.MarshalIndent(parser.GlobalImageConfig, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return "", err
	}

	if err = os.WriteFile(workDir+"/image_config.json", []byte(jsonData), 0644); err != nil {
		log.Println("Error writing file:", err)
		return "", err
	}

	sha256Value, err := common.CalculateFileSHA256(workDir + "/image_config.json")
	if err != nil {
		log.Println("Failed to calculate sha256")
		return "", err
	}

	sha256sums := fmt.Sprintf("%x", sha256Value)

	log.Println("Rename image json...")
	if err := os.Rename(workDir+"/image_config.json", workDir+"/"+sha256sums+".json"); err != nil {
		log.Printf("Error renaming folder: %v\n", err)
		return "", err
	}

	return sha256sums, nil
}

func ecsSaveCreateRepositories(workDir string, layerSha256 string, imageName string, imageTag string) error {
	result := fmt.Sprintf("{\n\t\"%s\":{\n\t\t\"%s\":\"%s\"\n\t}\n}", imageName, imageTag, layerSha256)

	if err := common.CreateFile(workDir+"/repositories", result); err != nil {
		return err
	}

	return nil
}

func ecsSaveCreateManifest(workDir string, layerSha256 string, configSha256 string, imageName string, imageTag string) error {

	parser.GlobalImageInfos[0].Config = configSha256 + ".json"
	parser.GlobalImageInfos[0].RepoTags[0] = imageTag
	parser.GlobalImageInfos[0].Layers[0] = layerSha256 + "/layer.tar"

	jsonData, err := json.MarshalIndent(parser.GlobalImageInfos, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return err
	}

	if err = os.WriteFile(workDir+"/manifest.json", []byte(jsonData), 0644); err != nil {
		log.Println("Error writing file:", err)
		return err
	}

	return nil
}

func EcsSave(bundlePath string, imageName string) error {

	_, err := os.Stat(bundlePath)
	if os.IsNotExist(err) {
		log.Printf("The bundle %s does not exist\n", bundlePath)
		return err
	}

	if err := ecsSaveCreateWorkDir("./ecs_temp"); err != nil {
		log.Println("Failed to create working directory!")
	}

	layerSha256, err := ecsSaveCreateLayer(bundlePath, "./ecs_temp")
	if err != nil {
		log.Println("Failed to create layer!")
	}

	configSha256, err := ecsSaveImageConfig(bundlePath, "./ecs_temp", layerSha256)
	if err != nil {
		log.Println("Failed to create layer!")
	}

	if err := ecsSaveCreateRepositories("./ecs_temp", layerSha256, imageName, "latest"); err != nil {
		log.Println("Failed to create repositories!")
	}

	if err := ecsSaveCreateManifest("./ecs_temp", layerSha256, configSha256, imageName, "latest"); err != nil {
		log.Println("Failed to create repositories!")
	}

	if err := common.TarDirectory(imageName, "./ecs_temp"); err != nil {
		log.Println("Failed to create image tar")
		return err
	}

	err = os.RemoveAll("./ecs_temp")
	if err != nil {
		log.Fatalf("The wrokDir removed faild, please retry.\n")
	}

	log.Println("create container tar sucess!")
	return nil
}
