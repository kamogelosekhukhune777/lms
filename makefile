# ============================================================================================
LMS_APP          := lms
BASE_IMAGE_NAME  := localhost/kamogelosekhukhune777
VERSION          := 0.0.1
VENDLY_IMAGE     := $(BASE_IMAGE_NAME)/$(LMS_APP):$(VERSION)

# ============================================================================================
# Building containers

build: lms

lms:
	docker build \
		-f zarf/docker/dockerfile.lms \
		-t $(LMS_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ============================================================================================
# Docker Compose

compose-up:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml -p compose up -d

compose-build-up: build compose-up

compose-down:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml down

compose-logs:
	cd ./zarf/compose/ && docker compose -f docker_compose.yaml logs

# ============================================================================================

run:
	go run api/services/lms/main.go | go run api/tooling/logfmt/main.go

version:
	go run api/services/lms/main.go --version

run-help:
	go run api/services/lms/main.go --help | go run api/tooling/logfmt/main.go

# ============================================================================================
# Metrics and Tracing

metrics-view-sc:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

metrics-view:
	expvarmon -ports="localhost:4020" -endpoint="/metrics" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open http://localhost:3010/debug/statsviz

# ============================================================================================
# Modules supports

tidy:
	go mod tidy
	go mod vendor