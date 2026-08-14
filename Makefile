MAKE_PREFIX := \033[1;34m[MAKE]\033[0m
log = @printf "$(MAKE_PREFIX) %s\n" "$(1)"

generate: mock wire tidy fmt

tidy:
	$(call log,Tidying Go modules...)

	go mod tidy

	$(call log,Go modules tidying complete.)

wire:
	$(call log,Running wire dependency injection...)

	wiresetgen generate

	wire ./cmd/...

	go fmt ./cmd/...

	$(call log,Wire generation complete.)

mock: mock-clean
	mockery

mock-clean:
	find . -type d -name "mock" -exec rm -rf {} +

fmt:
	$(call log,Formatting code...)

	go fmt ./...

	$(call log,Code formatting complete.)

start:
	$(call log,Starting the application...)

	go run ./cmd/server/main.go

start-cron:
	$(call log,Starting the cron processor...)

	go run ./cmd/cron/main.go -name=$(name)

	$(call log,Cron processor started.)

build:
	$(call log,Building the application...)

	go build -o bin/ ./cmd/...

	$(call log,Application build complete.)

tools-install:
	go install github.com/vektra/mockery/v3@v3.5.5