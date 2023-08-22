package builder

import (
	"docker-deployment/service"
	"log"
	"os"
	"os/exec"
)

func DockerImageBuilder(values Values, err error, directoryPath string) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	cmdBuild := exec.Command("docker", "build", "-t", values.ImageName+":"+values.ImageTag, "-f", values.DockerfilePath, ".")

	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err = cmdBuild.Run()
	if err != nil {
		log.Fatalf("Docker build command failed: %s", err)
	}

	filePath := directoryPath + "/" + values.ImageName + ".tar.gz"
	cmdSave := exec.Command("sh", "-c", "docker save "+values.ImageName+":"+values.ImageTag+" | gzip > "+filePath)
	cmdSave.Stdout = os.Stdout
	cmdSave.Stderr = os.Stderr
	err = cmdSave.Run()
	if err != nil {
		log.Fatalf("Docker save command failed: %s", err)
	}

	config := service.FileUploadConfig{
		FilePath: filePath,
		FieldValues: []service.Field{
			{"fileName", values.ImageName + ".tar"},
			{"dirName", "deployment/" + values.ImageName + "/" + values.ImageTag},
		},
		UploadURL:     values.UploadUrl,
		Authorization: values.Authorization,
	}
	err = service.PostFile(config)
	if err != nil {
		log.Fatalf("Post docker image failed %s - %s", filePath, err)
	}
	return err
}
