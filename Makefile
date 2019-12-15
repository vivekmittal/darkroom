APP=darkroom
APP_EXECUTABLE="./out/$(APP)"

all: test-ci

setup:
	go get golang.org/x/lint/golint

compile:
	mkdir -p out
	go build -o $(APP_EXECUTABLE) main.go

lint:
	@golint ./... | { grep -vwE "exported (var|function|method|type|const) \S+ should have comment" || true; }

format:
	go fmt ./...

vet:
	go vet ./...

test:
	go test ./... -covermode=count -coverprofile=profile.cov

copy-config:
	cp config.example.yaml config.yaml

docker-image:
	docker build -t ${USER}/darkroom:latest -f build/Dockerfile .

docker-docs:
	docker build -t darkroom-docs:latest -f build/Dockerfile.docs .

test-ci: copy-config compile lint format vet test
