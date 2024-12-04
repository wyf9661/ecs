package main

import (
	"archive/tar"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func tarDirectory(tarFileName string, sourceDir string) error {

	tarFile, err := os.Create(tarFileName)
	if err != nil {
		return err
	}
	defer tarFile.Close()

	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// Walk through the source directory
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		header.Name, err = filepath.Rel(sourceDir, filePath)
		if err != nil {
			return err
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tarWriter, file)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func calculateFileSHA256(filePath string) ([]byte, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}

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

	if err = createDir(workDir); err != nil {
		log.Println("Failed to create working directory")
		return err
	}

	return nil
}

func ecsSaveCreateLayerTar(bundlePath string, layImage string) error {
	if err := tarDirectory(layImage, bundlePath); err != nil {
		log.Println("Failed to create lay tar")
		return err
	}
	return nil
}

func ecsSaveCreateLayerJson(layerPath string, sha256Value string) error {

	GlobalLayerConfig.Id = sha256Value

	now := time.Now()
	formattedTime := now.Format("2006-01-02T15:04:05-07:00")
	GlobalLayerConfig.Created = formattedTime

	jsonData, err := json.MarshalIndent(GlobalLayerConfig, "", "    ")
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
	if err := createFile(layerPath+"/VERSION", "1.0"); err != nil {
		return err
	}
	return nil
}

func ecsSaveCreateLayer(bundlePath string, workDir string) (string, error) {
	log.Println("Creating layer...")

	log.Println("Creating layer directory...")
	if err := createDir(workDir + "/latest"); err != nil {
		log.Println("Failed to create working directory")
		return "", err
	}

	log.Println("Creating layer tarball...")
	if err := ecsSaveCreateLayerTar(bundlePath, workDir+"/latest"+"/layer.tar"); err != nil {
		log.Println("Failed to create working directory")
		return "", err
	}

	sha256Value, err := calculateFileSHA256(workDir + "/latest" + "/layer.tar")
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
	GlobalImageConfig.Created = formattedTime

	GlobalImageConfig.Architecture = GlobalConfigStruct.Platform.Arch
	GlobalImageConfig.Os = GlobalConfigStruct.Platform.Os
	GlobalImageConfig.Config.User = strconv.Itoa(GlobalConfigStruct.Process.User.Uid) + ":" + strconv.Itoa(GlobalConfigStruct.Process.User.Gid)
	GlobalImageConfig.Config.Env = GlobalConfigStruct.Process.Env
	GlobalImageConfig.Config.Entrypoint = GlobalConfigStruct.Process.Args
	GlobalImageConfig.Config.WorkingDir = GlobalConfigStruct.Process.Cwd

	GlobalImageConfig.Config.Labels["hostname"] = GlobalConfigStruct.Hostname

	for i := 0; i < len(GlobalConfigStruct.Mounts); i++ {
		GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".destination"] = GlobalConfigStruct.Mounts[i].Destination
		GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".source"] = GlobalConfigStruct.Mounts[i].Source
		GlobalImageConfig.Config.Labels["mounts."+strconv.Itoa(i)+".options"] = GlobalConfigStruct.Mounts[i].Options[0]
	}

	for i := 0; i < len(GlobalConfigStruct.Sylixos.Commands); i++ {
		GlobalImageConfig.Config.Labels["sylixos.commands."+strconv.Itoa(i)] = GlobalConfigStruct.Sylixos.Commands[i]
	}

	for i := 0; i < len(GlobalConfigStruct.Sylixos.Devices); i++ {
		GlobalImageConfig.Config.Labels["sylixos.devices."+strconv.Itoa(i)+".path"] = GlobalConfigStruct.Sylixos.Devices[i].Path
		GlobalImageConfig.Config.Labels["sylixos.devices."+strconv.Itoa(i)+".access"] = GlobalConfigStruct.Sylixos.Devices[i].Access
	}

	GlobalImageConfig.Config.Labels["sylixos.resources.cpu.highestPrio"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.Cpu.HighestPrio)
	GlobalImageConfig.Config.Labels["sylixos.resources.cpu.lowestPrio"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.Cpu.LowestPrio)

	GlobalImageConfig.Config.Labels["sylixos.resources.memory.kheapLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.Memory.KheapLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.memory.memoryLimitMB"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.Memory.MemoryLimitMB)

	GlobalImageConfig.Config.Labels["sylixos.resources.disk.limitMB"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.Disk.LimitMB)

	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.deviceLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.DeviceLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.dlopenLibraryLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.DlopenLibraryLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.eventLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.EventLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.eventSetLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.EventSetLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.msgQueueLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.MsgQueueLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.partitionLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.PartitionLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.posixMqueueLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.PosixMqueueLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.regionLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.RegionLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.rmsLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.RmsLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.socketLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.SocketLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.srtpLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.SrtpLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadPoolLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadPoolLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.threadVarLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.ThreadVarLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.timerLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.TimerLimit)
	GlobalImageConfig.Config.Labels["sylixos.resources.kernelObject.xsiipcLimit"] = strconv.Itoa(GlobalConfigStruct.Sylixos.Resources.KernelObject.XsiipcLimit)

	GlobalImageConfig.Config.Labels["sylixos.network.ftpdEnable"] = fmt.Sprintf("%t", GlobalConfigStruct.Sylixos.Network.FtpdEnable)
	GlobalImageConfig.Config.Labels["sylixos.network.telnetdEnable"] = fmt.Sprintf("%t", GlobalConfigStruct.Sylixos.Network.TelnetdEnable)
	GlobalImageConfig.Config.Labels["sylixos.network.sshdEnable"] = fmt.Sprintf("%t", GlobalConfigStruct.Sylixos.Network.SshdEnable)

	GlobalImageConfig.Rootfs.Type = "layers"
	GlobalImageConfig.Rootfs.DiffIds[0] = "sha256:" + layerSha256

	GlobalImageConfig.History[0].Created = formattedTime
	GlobalImageConfig.History[0].CreatedBy = "ecs (go version for built)"
	GlobalImageConfig.History[0].EmptyLayer = false

	jsonData, err := json.MarshalIndent(GlobalImageConfig, "", "    ")
	if err != nil {
		fmt.Println("Error parsing jsonData:", err)
		return "", err
	}

	if err = os.WriteFile(workDir+"/image_config.json", []byte(jsonData), 0644); err != nil {
		log.Println("Error writing file:", err)
		return "", err
	}

	sha256Value, err := calculateFileSHA256(workDir + "/image_config.json")
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

	if err := createFile(workDir+"/repositories", result); err != nil {
		return err
	}

	return nil
}

func ecsSaveCreateManifest(workDir string, layerSha256 string, configSha256 string, imageName string, imageTag string) error {

	GlobalImageInfos[0].Config = configSha256 + ".json"
	GlobalImageInfos[0].RepoTags[0] = imageTag
	GlobalImageInfos[0].Layers[0] = layerSha256 + "/layer.tar"

	jsonData, err := json.MarshalIndent(GlobalImageInfos, "", "    ")
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

func ecsSave(bundlePath string, imageName string) error {

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

	if err := tarDirectory(imageName, "./ecs_temp"); err != nil {
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
