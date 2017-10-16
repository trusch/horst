docker: textsplitter
	docker build -t trusch/horst/textsplitter .

textsplitter: $(shell find $(GOPATH)/src/github.com/trusch/horst -name "*.go")
	docker run \
		-v $(GOPATH):/go \
		-w /go/src/github.com/trusch/horst/cmd/textsplitter \
		-u $(shell stat -c '%u:%g' .) \
		-e "CGO_ENABLED=0" \
		-e "GOOS=linux" \
		golang:1.9 go build -v --ldflags '-extldflags "-static"' -o textsplitter
