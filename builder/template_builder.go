package builder

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"text/template"
)

type TemplateData struct {
	IMAGE_NAME                string
	IMAGE_TAG                 string
	BIND_CONTAINER_PORT       []string
	BIND_VOLUMES              []string
	MEMORY_LIMIT              string
	CPU_LIMIT                 string
	CONTAINER_NAME            string
	VOLUMES                   []string
	HEALTH_CHECK_URL          string
	HEALTH_CHECK_INTERVAL     string
	HEALTH_CHECK_TIMEOUT      string
	HEALTH_CHECK_RETRIES      int
	HEALTH_CHECK_START_PERIOD string
}

func TemplateBuilder(values Values, directoryPath string) error {

	data := &TemplateData{
		IMAGE_NAME:                values.ImageName,
		IMAGE_TAG:                 values.ImageTag,
		BIND_CONTAINER_PORT:       values.Bind.Ports,
		BIND_VOLUMES:              values.Bind.Volumes,
		MEMORY_LIMIT:              values.Limit.Memory,
		CPU_LIMIT:                 values.Limit.Cpu,
		VOLUMES:                   values.Volumes,
		CONTAINER_NAME:            values.ContainerName,
		HEALTH_CHECK_URL:          values.HealthCheck.Url,
		HEALTH_CHECK_INTERVAL:     values.HealthCheck.Interval,
		HEALTH_CHECK_TIMEOUT:      values.HealthCheck.Timeout,
		HEALTH_CHECK_RETRIES:      values.HealthCheck.Retries,
		HEALTH_CHECK_START_PERIOD: values.HealthCheck.StartPeriod,
	}

	if values.PullImageHost != "" {
		data.IMAGE_NAME = values.PullImageHost + "/" + data.IMAGE_NAME
	} else if values.PushImageHost != "" {
		data.IMAGE_NAME = values.PushImageHost + "/" + data.IMAGE_NAME
	}

	tmpl, err := template.ParseFiles(values.TemplatePath)
	if err != nil {
		log.Fatalf("Failed to load template file %s - %s", values.TemplatePath, err)
	}

	file, err := os.Create(directoryPath + "docker-compose.yml")

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Failed to close file: %w", err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}

	deployContent, err := LoadFileContentToBase64(directoryPath + "docker-compose.yml")

	deployRef := fmt.Sprintf("-appName \"%s\" -imageTag \"%s\" -deployment \"%s\"",
		values.ContainerName, values.ImageTag, deployContent)

	return WriteStringToFile("deploy-refs", deployRef)
}

func LoadFileContentToBase64(filePath string) (string, error) {
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileContent), nil
}

func WriteStringToFile(filename string, data string) error {
	err := os.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		log.Fatalf("Failed to write to file %s: %v", filename, err)
		return err
	}
	return nil
}
