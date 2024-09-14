docker-build:
	git submodule update --recursive --init
	docker-compose up --build -d
	docker-compose down

test:
	make docker-build
	./ci/deployment -imageTag "v5"  -config deployment.yaml -pushImage=false -storeImageAsFile=true