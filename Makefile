# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

include Makefile.common

GO     ?= GO15VENDOREXPERIMENT=1 go
GOARCH := $(shell $(GO) env GOARCH)
GOHOSTARCH := $(shell $(GO) env GOHOSTARCH)

PROMTOOL    ?= $(FIRST_GOPATH)/bin/promtool

DOCKER_IMAGE_NAME       ?= iscdhcpd-exporter
MACH                    ?= $(shell uname -m)
DOCKERFILE              ?= Dockerfile

STATICCHECK_IGNORE =

ifeq ($(OS),Windows_NT)
    OS_detected := Windows
else
    OS_detected := $(shell uname -s)
endif

ifeq ($(GOHOSTARCH),amd64)
	ifeq ($(OS_detected),$(filter $(OS_detected),Linux FreeBSD Darwin Windows))
                # Only supported on amd64
                test-flags := -race
        endif
endif


# By default, "cross" test with ourselves to cover unknown pairings.
$(eval $(call goarch_pair,amd64,386))
$(eval $(call goarch_pair,mips64,mips))
$(eval $(call goarch_pair,mips64el,mipsel))

all: style vet staticcheck build test

.PHONY: checkmetrics
checkmetrics: $(PROMTOOL)
	@echo ">> checking metrics for correctness"
	./checkmetrics.sh $(PROMTOOL) $(e2e-out)

.PHONY: docker
docker:
ifeq ($(MACH), ppc64le)
	$(eval DOCKERFILE=Dockerfile.ppc64le)
endif
	@echo ">> building docker image from $(DOCKERFILE)"
	@docker build --file $(DOCKERFILE) -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

.PHONY: test-docker
test-docker:
	@echo ">> testing docker image"
	./test_image.sh "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" 9367

.PHONY: promtool $(FIRST_GOPATH)/bin/promtool
$(FIRST_GOPATH)/bin/promtool promtool:
	@GOOS= GOARCH= $(GO) get -u github.com/prometheus/prometheus/cmd/promtool
