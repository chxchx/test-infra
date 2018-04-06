# Copyright Istio Authors 2018.

PROJECT ?= istio-testing
IMG := gcr.io/$(PROJECT)/kettle
TAG ?= 0.1.0

build:
	docker build -t $(IMG):$(TAG) .

push: build
	gcloud docker -- push $(IMG):$(TAG)

run:
	docker run $(IMG):$(TAG)

run.bash:
	docker run -it --entrypoint bash --privileged $(IMG):$(TAG)

.PHONY: build push run run.bash
