package builder

import (
	"flag"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Bind struct {
	Ports   []string `yaml:"ports"`
	Volumes []string `yaml:"volumes"`
}

type Limit struct {
	Memory string `yaml:"memory"`
	Cpu    string `yaml:"cpu"`
}

type Healthcheck struct {
	Url         string `yaml:"url"`
	Interval    string `yaml:"interval"`
	Timeout     string `yaml:"timeout"`
	Retries     int    `yaml:"retries"`
	StartPeriod string `yaml:"start_period"`
}

type Values struct {
	ImageName      string            `yaml:"imageName"`
	ContainerName  string            `yaml:"containerName"`
	ImageTag       string            `yaml:"imageTag"`
	TemplatePath   string            `yaml:"template"`
	DockerfilePath string            `yaml:"dockerfile"`
	Bind           Bind              `yaml:"bind"`
	Limit          Limit             `yaml:"limit"`
	Volumes        []string          `yaml:"volumes"`
	Environment    map[string]string `yaml:"environment"`
	Args           map[string]string `yaml:"args"`
	UploadUrl      string            `yaml:"uploadUrl"`
	PushImage      bool              `yaml:"pushImage"`
	PushImageHost  string            `yaml:"pushImageHost"`
	PullImageHost  string            `yaml:"pullImageHost"`
	LoginUsername  string            `yaml:"loginUsername"`
	LoginPassword  string            `yaml:"loginPassword"`
	HealthCheck    Healthcheck       `yaml:"healthcheck"`
}

func ValuesBuilder() Values {
	var config Values
	loginUsername := flag.String("loginUsername", "", "Docker login username, if provided username and password, will try to login on Docker.")
	loginPassword := flag.String("loginPassword", "", "Docker login password, if provided username and password, will try to login on Docker.")
	configPath := flag.String("config", "cd/deployment.yaml", "Path to the deployment.yaml file.")
	pushImage := flag.Bool("pushImage", true, "A flag to indicate whether to push the image or not. If true the generate docker image and docker-compose.yaml, will be pushed.")
	imageTag := flag.String("imageTag", "", "The imageTag parameter is used during the Docker image build process to tag the image that is being built.")

	flag.Parse()
	log.Println("PushImage === ", *pushImage)

	if *configPath != "" {
		file, err := os.ReadFile(*configPath)

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
			config.TemplatePath = "./ci/deployment.mustache"
		}

		config.PushImage = *pushImage
		config.LoginPassword = *loginPassword
		config.LoginUsername = *loginUsername

		if *imageTag != "" {
			config.ImageTag = *imageTag
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

	} else {
		log.Fatalf("Path to the deployment.yaml file and authorization  must be provided.")
	}
	return config
}
