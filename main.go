package main

import (
	"docker-deployment/builder"
	"log"
	"os"
)

func main() {
	values := builder.ValuesBuilder()

	directoryPath := ".temp/"

	err := os.MkdirAll(directoryPath, os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to create directory: %s", err)
	}

	err = builder.DockerImageBuilder(values, err, directoryPath)

	builder.TemplateBuilder(values, directoryPath)
}
