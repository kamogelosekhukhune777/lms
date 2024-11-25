
# ============================================================================================

run:
	go run api/services/vendly/main.go | go run api/tooling/logfmt/main.go

version:
	go run api/services/vendly/main.go --version

run-help:
	go run api/services/vendly/main.go --help | go run api/tooling/logfmt/main.go

# ============================================================================================
# Modules supports

tidy:
	go mod tidy
	go mod vendor