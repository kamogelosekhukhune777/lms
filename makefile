
# ============================================================================================

run:
	go run api/services/vendly/main.go | go run api/tooling/logfmt/main.go

version:
	go run api/services/vendly/main.go --version

run-help:
	go run api/services/vendly/main.go --help | go run api/tooling/logfmt/main.go

# ==============================================================================
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