package builder

import (
	"docker-deployment/service"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func DockerImageBuilder(values Values, err error, directoryPath string) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	//cmdBuild := exec.Command("docker", "build", "-t", values.ImageName+":"+values.ImageTag, "-f", values.DockerfilePath, ".")
	//
	//cmdBuild.Stdout = os.Stdout
	//cmdBuild.Stderr = os.Stderr
	//err = cmdBuild.Run()
	//if err != nil {
	//	log.Fatalf("Docker build command failed: %s", err)
	//}

	filePath := directoryPath + "/" + values.ImageName + ".tar.gz"
	//cmdSave := exec.Command("sh", "-c", "docker save "+values.ImageName+":"+values.ImageTag+" | gzip > "+filePath)
	//cmdSave.Stdout = os.Stdout
	//cmdSave.Stderr = os.Stderr
	//err = cmdSave.Run()
	//if err != nil {
	//	log.Fatalf("Docker save command failed: %s", err)
	//}

	maxSize := 75
	if !isAGoodFileSize(filePath, int64(maxSize)) {
		err = splitFile(values, strconv.Itoa(maxSize)+"M", err, directoryPath, filePath)
	} else {
		config := service.FileUploadConfig{
			FilePath: filePath,
			FieldValues: []service.Field{
				{"fileName", values.ImageName + ".tar"},
				{"dirName", "deployment/" + values.ImageName + "/" + values.ImageTag},
			},
			UploadURL:     values.UploadUrl,
			Authorization: values.Authorization,
		}
		err = PostFile(err, config, filePath)
	}

	return err
}

func isAGoodFileSize(filePath string, goodSize int64) bool {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Fatal(err)
	}
	fileSizeInBytes := fileInfo.Size()
	fileSizeInMB := fileSizeInBytes / (1024 * 1024)
	return fileSizeInMB < goodSize
}

func splitFile(values Values, maxSize string, err error, directoryPath string, filePath string) error {
	cmdSplit := exec.Command("split", "-b", maxSize, filePath, directoryPath+"/"+values.ImageName+".part-")
	cmdSplit.Stdout = os.Stdout
	cmdSplit.Stderr = os.Stderr
	err = cmdSplit.Run()
	if err != nil {
		log.Fatalf("Failed to split file: %s", err)
	}

	uploadFilePart(values, directoryPath)
	return err
}

func uploadFilePart(values Values, directoryPath string) {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePath := directoryPath + filepath.Base(path)

			if isAGoodFileSize(filePath, 76) {
				fileName := strings.Split(filepath.Base(path), ".part-")[0]

				config := service.FileUploadConfig{
					FilePath: filePath,
					FieldValues: []service.Field{
						{"fileName", fileName},
						{"dirName", "deployment/" + values.ImageName + "/" + values.ImageTag},
					},
					UploadURL:     values.UploadUrl,
					Authorization: values.Authorization,
				}
				err = PostFile(err, config, filePath)
			} else {
				log.Printf("File %s too large, than will be ignored", filePath)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %v: %v\n", directoryPath, err)
	}
}

func PostFile(err error, config service.FileUploadConfig, filePath string) error {
	err = service.PostFile(config)
	if err != nil {
		log.Fatalf("Post docker image failed %s - %s", filePath, err)
	}
	return err
}
