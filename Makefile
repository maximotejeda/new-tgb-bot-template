# must create a .env file with info
# must have compose installed
include .env
export
OS:=${shell go env GOOS}
ARCH=$(shell go env GOARCH)
OOSS="linux"
ARRCHS="arm 386"
DEBUG=1
SERVICE=new-tgb-bot-template
VERSION=0.0.0_1
BINAME=$(SERVICE)-$(OS)-$(ARCH)-$(VERSION)
BINAMEARM=$(SERVICE)-$(OS)-arm64-$(VERSION)
# can be docker or podman or whatever
CONTAINERS=docker
COMPOSE=$(CONTAINERS)-compose
# Configure local registry
REGADDR=192.168.0.151:32000
K8SRSNAME=$(shell kubectl get rs --no-headers -o custom-columns=":metadata.name" | grep us-dop-bot)
.phony: all clean build test clean-image build-image build-image-debug run-image run-image-debug run-local


build-image: build 
# here we made the images and push to registry with buildx 
	@$(CONTAINERS) buildx build --build-arg="BINAME=$(BINAMEARM)" --platform linux/arm64 --push -t $(REGADDR)/$(SERVICE):latest .

# Here we upload it to local 

build-test-image:
	@$(CONTAINERS) buildx build --platform linux/arm64 --push  -t $(REGADDR)/$(SERVICE):latest .

run-image: build-image
	@$(CONTAINERS) compose -f docker-compose.yaml up

build-image-debug: clean
	@$(CONTAINERS) compose -f docker-compose-debug.yaml build

run-image-debug: build-image-debug
	@$(CONTAINERS) compose -f docker-compose-debug.yaml up

run-local:clean build
	@bin/$(BINAME)

build: clean
	#@mkdir dolardb
	@env GOOS=$(OS) GOARCH=$(arch) go build -o ./bin/$(BINAME) ./cmd/.
	@env GOOS=$(OS) GOARCH=arm64 go build -o ./bin/$(BINAMEARM) ./cmd/.

create-descriptors:
	@envsubst < k8s/deployment.yml.template > k8s/deployment.yml

deploy: build-image create-descriptors
	@kubectl apply -f k8s/deployment.yml
	@kubectl scale rs $(K8SRSNAME) --replicas=0
	@kubectl scale rs  $(K8SRSNAME) --replicas=1

test:
	@go -count=1 test ./...
clean:
	@rm -rf ./bin 

clean-image:
	@$(CONTAINERS) system prune -f


