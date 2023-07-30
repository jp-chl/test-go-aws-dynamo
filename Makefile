.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/hello src/main.go

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose

remove:
	serverless remove
