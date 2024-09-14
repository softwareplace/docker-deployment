docker-build:
	docker-compose up --build -d
	docker-compose down

test:
	make docker-build
	./ci/deployment \
		-imageTag "v5" \
		-config "deployment.yaml" \
		-fileName "./.temp/my-docker-image.tar" \
		-storeImageAsFile true \
		-pushImage false
