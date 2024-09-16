package builder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func DockerImageBuilder(values Values) error {
	if values.DockerfilePath == "." {
		values.DockerfilePath = "./Dockerfile"
	}

	err := dockerRun(values)
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

func EnvGenerate(values Values, directoryPath string) {
	imageNameTag := values.ImageName + ":" + values.ImageTag
	image := values.ImageName + ":" + values.ImageTag

	if values.PullImageHost != "" {
		image = values.PushImageHost + "/" + values.ImageName + ":" + values.ImageTag
	}

	template := "export %s=%s\n"
	var result = fmt.Sprintf(template, "DOCKER_IMAGE_NAME", imageNameTag)
	result += fmt.Sprintf(template, "DOCKER_IMAGE_FULL_NAME", image)
	result += fmt.Sprintf(template, "DOCKER_CONTAINER_NAME", values.ContainerName)

	err := os.WriteFile(directoryPath+".env", []byte(result), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Printf("Declared env:\n%sTo load theses envs, run\n\tsource %s\n", result, directoryPath+".env")

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
	fmt.Printf("Declared env DOCKER_IMAGE_FULL_NAME=%s", image)
	return err
}

func showDoLogin(username string, password string) bool {
	return username != "" && password != ""
}

func dockerRun(values Values) error {
	imageNameTag := values.ImageName + ":" + values.ImageTag

	log.Printf("Building doker image %s started", imageNameTag)

	args := []string{"build", "-t", imageNameTag, "-f", values.DockerfilePath}

	for argKey, argValue := range values.Args {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", argKey, argValue))
	}

	for _, host := range values.ExtrasHosts {
		args = append(args, "--add-host", host)
	}

	args = append(args, ".")

	cmdBuild := exec.Command("docker", args...)

	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err := cmdBuild.Run()

	if err != nil {
		log.Fatalf("Docker build command failed: %s", err)
	} else {
		if values.StoreImageAsFile {
			err = storeDockerImageToFile(values)
		}
	}
	return err
}

func storeDockerImageToFile(values Values) error {
	log.Println("Storing image as file")
	image := values.ImageName + ":" + values.ImageTag

	if values.PullImageHost != "" {
		image = values.PushImageHost + "/" + values.ImageName + ":" + values.ImageTag
	}

	err := rumCommand("docker", "image", "tag", values.ImageName+":"+values.ImageTag, image)

	if err != nil {
		return err
	}

	err = rumCommand("docker", "save", "-o", values.FileName, image)
	if err != nil {
		return err
	}
	return err
}
