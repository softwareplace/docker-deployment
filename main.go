package main

import (
	"docker-deployment/builder"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	values := builder.ValuesBuilder()

	directoryPath := ".temp/"

	err := os.MkdirAll(directoryPath, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}

	err = builder.TemplateBuilder(values, directoryPath)

	if err != nil {
		log.Fatal(err)
	}

	_ = builder.DockerImageBuilder(values)

	builder.EnvGenerate(values, directoryPath)

	fmt.Println(directoryPath+"docker-compose.yml", "generated successful")

	cmdBuild := exec.Command("cat", directoryPath+"docker-compose.yml")

	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	err = cmdBuild.Run()

}
