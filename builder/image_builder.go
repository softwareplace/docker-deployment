package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func DockerImageBuilder(values Values, err error) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	err = dockerRun(values, err)
	if err != nil {
		return err
	}

	if values.PushImage == true {
		log.Println("Push image", values.PushImage)
		return dockerImageStorage(values)
	} else {
		log.Println("Push image was disable.")
	}

	return nil
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
	image := values.ImageName + ":" + values.ImageTag

	if values.PullImageHost != "" {
		image = values.PushImageHost + "/" + values.ImageName + ":" + values.ImageTag
	}

	log.Printf("Pushing image: %s\n", image)
	err := rumCommand("docker", "image", "tag", values.ImageName+":"+values.ImageTag, image)
	if err != nil {
		return err
	}

	if showDoLogin(values.LoginUsername, values.LoginPassword) {
		err = rumCommand("docker", "login", "-p", values.LoginPassword, "-u", values.LoginUsername, values.PushImageHost)
		if err != nil {
			return err
		}
	}

	err = rumCommand("docker", "push", image)
	if err != nil {
		return err
	}

	if showDoLogin(values.LoginUsername, values.LoginPassword) {
		err = rumCommand("docker", "logout", values.PushImageHost)
		if err != nil {
			return err
		}
	}

	return err
}

func showDoLogin(username string, password string) bool {
	return username != "" && password != ""
}

func dockerRun(values Values, err error) error {
	log.Printf("Building doker image %s", values.ImageName+":"+values.ImageTag)

	args := []string{"build", "-t", values.ImageName + ":" + values.ImageTag, "-f", values.DockerfilePath}

	for argKey, argValue := range values.Args {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", argKey, argValue))
	}

	args = append(args, ".")

	cmdBuild := exec.Command("docker", args...)

	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err = cmdBuild.Run()
	if err != nil {
		log.Fatalf("Docker build command failed: %s", err)
	}

	return err
}
