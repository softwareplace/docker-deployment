# Continuous Integration deployment builder

> Requires
> - Docker installed


## Args

| Arg            | Required | Default            | Description                                                                                                                     |
|----------------|----------|--------------------|---------------------------------------------------------------------------------------------------------------------------------|
| -loginUsername | true     |                    | Docker login username, if provided username and password, will try to login on Docker.                                          |
| -loginPassword | true     |                    | Docker login password, if provided username and password, will try to login on Docker.                                          |
| -help          | false    |                    | Show available args.                                                                                                            |
| -config        | false    | cd/deployment.yaml | Path to the deployment.yaml file.                                                                                               |
| -pushImage     | false    | true               | A flag to indicate whether to push the image or not. If true the generate docker image and docker-compose.yaml, will be pushed. |
| -imageTag      | false    |                    | The imageTag parameter is used during the Docker image build process to tag the image that is being built.                      |

- Config build

```yaml
imageName: my-app
containerName: my-app
# By running /deployment specifying --imageTag, the imageTag on configuration file will be ignored
imageTag: v1
template: "ci/deployment.mustache"
# If you want to push the image to a custom docker registry
pushImageHost: docker-registry.com
# The host that image should be pulled
# If you want to pull the image from a docker registry
# Recommended when the registry is running on the same host. If not set, pushImageHost will be used
pullImageHost: localhost:5000
# Full path that docker will check container health
healthcheck:
  retries: 5
  timeout: 15s
  interval: 30s
  start_period: 10s
  url: "http://172.17.0.1:8080"
dockerfile: .
bind:
  ports:
    - "127.0.0.1:8080:80"
    - "172.17.0.1:8080:80"
  volumes:
    - ~/.storage/:/root/.storage/
limit:
  memory: 2G
  cpu: 0.5
volumes:
  - "example_logs:"
```

