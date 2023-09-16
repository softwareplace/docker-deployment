package builder

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"text/template"
)

type TemplateData struct {
	ImageName              string
	ImageTag               string
	BindContainerPort      []string
	BindVolumes            []string
	MemoryLimit            string
	CpuLimit               string
	ContainerName          string
	Volumes                []string
	HealthCheckUrl         string
	HealthCheckInterval    string
	HealthCheckTimeout     string
	HealthCheckRetries     int
	HealthCheckStartPeriod string
}

func TemplateBuilder(values Values, directoryPath string) error {

	data := &TemplateData{
		ImageName:              values.ImageName,
		ImageTag:               values.ImageTag,
		BindContainerPort:      values.Bind.Ports,
		BindVolumes:            values.Bind.Volumes,
		MemoryLimit:            values.Limit.Memory,
		CpuLimit:               values.Limit.Cpu,
		Volumes:                values.Volumes,
		ContainerName:          values.ContainerName,
		HealthCheckUrl:         values.HealthCheck.Url,
		HealthCheckInterval:    values.HealthCheck.Interval,
		HealthCheckTimeout:     values.HealthCheck.Timeout,
		HealthCheckRetries:     values.HealthCheck.Retries,
		HealthCheckStartPeriod: values.HealthCheck.StartPeriod,
	}

	if values.PullImageHost != "" {
		data.ImageName = values.PullImageHost + "/" + data.ImageName
	} else if values.PushImageHost != "" {
		data.ImageName = values.PushImageHost + "/" + data.ImageName
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

	deployRef := fmt.Sprintf("-appName \"%s\" -ImageTag \"%s\" -deployment \"%s\"",
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
