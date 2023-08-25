package builder

import (
	"log"
	"os"
	"os/exec"
)

func DockerImageBuilder(values Values, err error) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	err = dockerRun(values, err)

	log.Println("Push image", values.PushImage)
	if values.PushImage == true {
		err = dockerImageStorage(values)
	}

	if err != nil {
		return err
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
	err = rumCommand("docker", "login", "-p", values.LoginPassword, "-u", values.LoginUsername, values.PushImageHost)
	if err != nil {
		return err
	}
	err = rumCommand("docker", "push", image)
	if err != nil {
		return err
	}
	err = rumCommand("docker", "logout", values.PushImageHost)
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
