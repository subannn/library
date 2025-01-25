LIBRARY_BINARY=libraryService

up_build: build_library
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

build_library:
	@echo "Building auth library..."
	cd ../src && env GOOS=linux CGO_ENABLED=0 go build -o ${LIBRARY_BINARY} ../src/server
	@echo "Done!"


