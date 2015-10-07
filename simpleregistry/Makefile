cross-build: linux darwin

linux:
	 CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o bin/linux-amd64-simpleregistry .

darwin:
	 CGO_ENABLED=0 GOOS=darwin go build -a -tags netgo -ldflags '-w' -o bin/darwin-amd64-simpleregistry .

docker:
	docker build -t binder/registry .
