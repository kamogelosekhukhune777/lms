
LMS_APP         := lms
BASE_IMAGE_NAME := localhost/kamogelosekhukhune777
VERSION         := 0.0.1
LMS_IMAGE       := $(BASE_IMAGE_NAME)/$(LMS_APP):$(VERSION)

# ==========================================================================================
# Building containers

build: lms

lms:
	docker build \
		-f zarf/docker/dockerfile.lms \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==========================================================================================
# Docker Compose

# ==========================================================================================
# Metrics and Tracing

# ==========================================================================================
# Running tests within the local computer

# ==========================================================================================
# Hitting endpoints

# ==========================================================================================
# shortcuts 

run:
	go run api/services/lms/main.go | go run api/tooling/logfmt/main.go

# ==========================================================================================
# Modules support

# ==========================================================================================
# Admin Frontend

# ==========================================================================================
# Help command