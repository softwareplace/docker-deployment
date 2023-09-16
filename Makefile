docker-build:
	git submodule update --recursive --init
	docker-compose up --build