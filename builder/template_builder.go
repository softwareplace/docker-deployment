package builder

import (
	"docker-deployment/service"
	"log"
	"os"
	"text/template"
)

type TemplateData struct {
	IMAGE_NAME          string
	IMAGE_TAG           string
	BIND_CONTAINER_PORT []string
	BIND_VOLUMES        []string
	MEMORY_LIMIT        string
	CPU_LIMIT           string
	VOLUMES             []string
}

func TemplateBuilder(values Values, directoryPath string) {

	data := &TemplateData{
		IMAGE_NAME:          values.ImageName,
		IMAGE_TAG:           values.ImageTag,
		BIND_CONTAINER_PORT: values.Bind.Ports,
		BIND_VOLUMES:        values.Bind.Volumes,
		MEMORY_LIMIT:        values.Limit.Memory,
		CPU_LIMIT:           values.Limit.Cpu,
		VOLUMES:             values.Volumes,
	}

	tmpl, err := template.ParseFiles(values.TemplatePath)
	if err != nil {
		log.Fatalf("Failed to load template file %s - %s", values.TemplatePath, err)
	}

	file, err := os.Create(directoryPath + "/docker-compose.yml")

	if err != nil {
		log.Fatal(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	err = tmpl.Execute(file, data)
	if err != nil {
		log.Fatal(err)
	}

	config := service.FileUploadConfig{
		FilePath: directoryPath + "/docker-compose.yml",
		FieldValues: []service.Field{
			{"fileName", "docker-compose"},
			{"dirName", "deployment/" + values.ImageName + "/" + values.ImageTag},
		},
		UploadURL:     values.UploadUrl,
		Authorization: values.Authorization,
	}
	err = service.PostFile(config)
	if err != nil {
		log.Fatalf("%s/docker-compose.yml - %s", directoryPath, err)
	}
}
