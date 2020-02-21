VERSION=0.0.1
TAG=vlzhvlzh/editor-js-link-resolver

.PHONY: all build publish

all: build publish

build:
	docker build -t $(TAG):$(VERSION) -t $(TAG):latest .

publish:
	docker push $(TAG):$(VERSION)
	docker push $(TAG):latest

