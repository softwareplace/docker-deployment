package builder

import (
	"docker-deployment/service"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var maxFileSize = 10

func DockerImageBuilder(values Values, err error) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	err = dockerRun(values, err)

	if values.PushImage {
		err = dockerImageStorage(values)
	}

	return err
}

func rumCommand(name string, arg ...string) error {
	cmdSave := exec.Command(name, arg...)
	cmdSave.Stdout = os.Stdout
	cmdSave.Stderr = os.Stderr
	err := cmdSave.Run()
	if err != nil {
		log.Fatalf("Running command %s failed: %s", name, err)
	}
	return err
}

func dockerImageStorage(values Values) error {
	image := values.PushImageHost + "/" + values.ImageName + ":" + values.ImageTag

	log.Printf("Pushing image: %s\n", image)
	err := rumCommand("docker", "image", "tag", values.ImageName+":"+values.ImageTag, image)
	if err != nil {
		return err
	}
	err = rumCommand("docker", "push", image)
	if err != nil {
		return err
	}
	return err
}

func dockerRun(values Values, err error) error {
	log.Printf("Building doker image %s", values.ImageName+":"+values.ImageTag)

	cmdBuild := exec.Command("docker", "build", "-t", values.ImageName+":"+values.ImageTag, "-f", values.DockerfilePath, ".")

	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err = cmdBuild.Run()
	if err != nil {
		log.Fatalf("Docker build command failed: %s", err)
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
	err = os.MkdirAll(directoryPath+"/parts", os.ModePerm)

	cmdSplit := exec.Command("split", "-b", maxSize, filePath, directoryPath+"/parts/"+values.ImageName+".part-")

	cmdSplit.Stdout = os.Stdout
	cmdSplit.Stderr = os.Stderr
	err = cmdSplit.Run()
	if err != nil {
		log.Fatalf("Failed to split file: %s", err)
	}

	uploadFilePart(values, directoryPath+"parts/")
	return err
}

func uploadFilePart(values Values, directoryPath string) {
	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			filePath := directoryPath + filepath.Base(path)

			if isAGoodFileSize(filePath, int64(maxFileSize+1)) {
				config := service.FileUploadConfig{
					FilePath: filePath,
					FieldValues: []service.Field{
						{"fileName", values.ImageName},
						{"dirName", "deployment/" + values.ImageName + "/" + values.ImageTag + "/parts"},
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
