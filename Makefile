BIN_DIR = $(PWD)/bin

# Will setup linting tools
setup-linting:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin v1.46.1
	chmod +x ./bin/golangci-lint

lint-docker:
	docker run --rm -i hadolint/hadolint < Dockerfile

# Download dependencies
install:
	go mod download

# remove unused dependencies
tidy:
	go mod tidy -v

# Runs project
run:
	go run app/cmd/server/main.go

test:
	go test ./...

test-coverage:
	go test -tags testing -v -cover -covermode=atomic ./...

fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

clean:
	if [ -f ${BIN_DIR} ] ; then rm ${BIN_DIR} ; fi

lint:
	./bin/golangci-lint run ./...

build:
	@echo "Building application"
	go build -o $(BIN_DIR)/robinlb app/cmd/server/main.go

all: install lint test