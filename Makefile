DOCKER_PREFIX=trusch/horst/

COMPONENTS = logger projector textfilter texthistogram textsanitizer textsplitter twittersource
BINARIES = $(foreach component,$(COMPONENTS),bin/$(component)) bin/horst-dev-runner bin/horst-k8s-compiler
DOCKER_IMAGES = $(foreach component,$(COMPONENTS),docker/$(component)) docker/horst-dev-runner docker/horst-k8s-compiler

docker-images: $(DOCKER_IMAGES)

binaries: $(BINARIES)

install: binaries
	cp bin/* $(GOPATH)/bin/

clean:
	rm -rf bin $(foreach component,$(COMPONENTS),docker/$(component)) docker/horst-dev-runner docker/horst-k8s-compiler vendor

bin/horst-dev-runner: $(shell find ./cmd/horst-dev-runner ./config) vendor
	mkdir -p bin
	docker run --rm -it \
		-u $(shell stat -c '%u:%g' .) \
		-v $(shell pwd):/go/src/github.com/trusch/horst \
		-w /go/src/github.com/trusch/horst \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		golang:1.9 go build -a -v -ldflags '-extldflags "-static"' -o bin/horst-dev-runner github.com/trusch/horst/cmd/horst-dev-runner

bin/horst-k8s-compiler: $(shell find ./cmd/horst-k8s-compiler ./config) vendor
	mkdir -p bin
	docker run --rm -it \
		-u $(shell stat -c '%u:%g' .) \
		-v $(shell pwd):/go/src/github.com/trusch/horst \
		-w /go/src/github.com/trusch/horst \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		golang:1.9 go build -a -v -ldflags '-extldflags "-static"' -o bin/horst-k8s-compiler github.com/trusch/horst/cmd/horst-k8s-compiler


bin/%: ./cmd/%/*.go ./components/%/*.go ./components/base/*.go ./components/interface.go ./runner/*.go ./config/*.go vendor
	mkdir -p bin
	docker run --rm -it \
		-u $(shell stat -c '%u:%g' .) \
		-v $(shell pwd):/go/src/github.com/trusch/horst \
		-w /go/src/github.com/trusch/horst \
		-e CGO_ENABLED=0 \
		-e GOOS=linux \
		golang:1.9 go build -a -v -ldflags '-extldflags "-static"' -o bin/$(shell basename $@) github.com/trusch/horst/cmd/$(shell basename $@)

docker/%: bin/% docker/Dockerfile.template
	cp $< docker/
	cd docker && sed s/{{binary}}/$(shell basename $@)/g Dockerfile.template > Dockerfile && docker build -t $(DOCKER_PREFIX)$(shell basename $@) .

bin:
	mkdir -p bin

vendor: Gopkg.toml
	docker run -it --rm \
		-v $(shell pwd):/go/src/github.com/trusch/horst \
		-w /go/src/github.com/trusch/horst \
		golang:1.9 sh -c "go get -v github.com/golang/dep/cmd/dep && dep ensure -v"
	docker run --rm -v $(shell pwd):/app -w /app busybox chown -R $(shell stat -c '%u:%g' .) vendor Gopkg.toml Gopkg.lock
