
DOCKER ?= $(shell which docker)
# These docs are in an order that determines how they show up in the PDF/HTML docs.
DOC_FILES := \
	README.md \
	code-of-conduct.md \
	project.md \
	media-types.md \
	inline-png-media-types.md \
	manifest.md \
	serialization.md

FIGURE_FILES := \
	media-types.png

default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'fmt' - format the json with indentation"
	@echo " * 'validate' - build the validation tool"

fmt:
	for i in *.json ; do jq --indent 2 -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done

docs: output/docs.pdf output/docs.html

output/docs.pdf: $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p output/ && \
	cp *.png $(shell pwd)/output && \
	$(DOCKER) run \
	-it \
	--rm \
	-v $(shell pwd)/:/input/:ro \
	-v $(shell pwd)/output/:/output/ \
	-u $(shell id -u) \
	vbatts/pandoc -f markdown_github -t latex -o /output/docs.pdf $(patsubst %,/input/%,$(DOC_FILES)) && \
	ls -sh $(shell readlink -f $@)

output/docs.html: $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p output/ && \
	cp *.png $(shell pwd)/output && \
	$(DOCKER) run \
	-it \
	--rm \
	-v $(shell pwd)/:/input/:ro \
	-v $(shell pwd)/output/:/output/ \
	-u $(shell id -u) \
	vbatts/pandoc -f markdown_github -t html5 -o /output/docs.html $(patsubst %,/input/%,$(DOC_FILES)) && \
	ls -sh $(shell readlink -f $@)

code-of-conduct.md:
	curl -o $@ https://raw.githubusercontent.com/opencontainers/tob/d2f9d68c1332870e40693fe077d311e0742bc73d/code-of-conduct.md

validate-examples: oci-validate-examples
	./oci-validate-examples < manifest.md

oci-validate-json: validate.go
	go build ./cmd/oci-validate-json

oci-validate-examples: cmd/oci-validate-examples/main.go
	go build ./cmd/oci-validate-examples

oci-image-tool:
	go build ./cmd/oci-image-tool

lint:
	for d in $(shell find . -type d -not -iwholename '*.git*'); do echo "$${d}" && ./lint "$${d}"; done

test:
	go test -race ./...

%.png: %.dot
	dot -Tpng $^ > $@

inline-png-%.md: %.png
	@printf '<img src="data:image/png;base64,%s" alt="$*"/>\n' "$(shell base64 $^)" > $@

.PHONY: \
	validate-examples \
	oci-image-tool \
	lint \
	docs \
	test
