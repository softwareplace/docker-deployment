package builder

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Bind struct {
	Ports   []string `yaml:"ports"`
	Volumes []string `yaml:"volumes"`
}

type Limit struct {
	Memory string `yaml:"memory"`
	Cpu    string `yaml:"cpu"`
}

type Values struct {
	ImageName      string   `yaml:"imageName"`
	ImageTag       string   `yaml:"imageTag"`
	TemplatePath   string   `yaml:"template"`
	DockerfilePath string   `yaml:"dockerfile"`
	Bind           Bind     `yaml:"bind"`
	Limit          Limit    `yaml:"limit"`
	Volumes        []string `yaml:"volumes"`
	Authorization  string   `yaml:"authorization"`
	UploadUrl      string   `yaml:"uploadUrl"`
}

func ValuesBuilder() Values {
	var config Values
	configPath := flag.String("config", "", "Path to the deployment.yaml file.")
	authorization := flag.String("authorization", "", "Api push docker config authorization token.")

	flag.Parse()

	if *configPath != "" && *authorization != "" {
		file, err := ioutil.ReadFile(*configPath)

		if err != nil {
			log.Fatalf("Error reading YAML file: %s\n", err)
		}

		err = yaml.Unmarshal(file, &config)

		if err != nil {
			log.Fatalf("Error parsing YAML file: %s\n", err)
		}

		// Check if mandatory values are set
		if config.ImageName == "" || len(config.Bind.Ports) == 0 || len(config.Bind.Volumes) == 0 {
			log.Fatalf("Error: Missing required parameters in YAML. ImageName, Ports and Volumes must be set")
		}
		// Set defaults
		if config.TemplatePath == "" {
			config.TemplatePath = "./docker-compose-template.yml"
		}

		if config.ImageTag == "" {
			config.ImageTag = "latest"
		}
		if config.DockerfilePath == "" {
			config.DockerfilePath = "./Dockerfile"
		}
		if config.Limit.Memory == "" {
			config.Limit.Memory = "512mb"
		}
		if config.Limit.Cpu == "" {
			config.Limit.Cpu = "0.5"
		}

		config.Authorization = *authorization
	} else {
		log.Fatalf("Path to the deployment.yaml file and authorization  must be provided.")
	}
	return config
}
