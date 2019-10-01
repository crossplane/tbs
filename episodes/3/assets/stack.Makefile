include stack.env

GO111MODULE ?= on
export GO111MODULE

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


STACK_PACKAGE=stack-package
STACK_PACKAGE_REGISTRY=$(STACK_PACKAGE)/.registry

all: build
.PHONY: all

clean: clean-stack-package clean-binary
.PHONY: clean

clean-stack-package:
	rm -r $(STACK_PACKAGE)
.PHONY: clean-stack-package

clean-binary:
	rm -r bin
.PHONY: clean-binary

build: bundle docker-build
.PHONY: build

publish: docker-push
.PHONY: publish

# Initialize the stack bundle folder
$(STACK_PACKAGE_REGISTRY):
	mkdir -p $(STACK_PACKAGE_REGISTRY)/resources
	touch $(STACK_PACKAGE_REGISTRY)/app.yaml $(STACK_PACKAGE_REGISTRY)/install.yaml

bundle: $(STACK_PACKAGE_REGISTRY)
	# Copy CRDs over
	#
	# The reason this looks complicated is because it is
	# preserving the original crd filenames and changing
	# *.yaml to *.crd.yaml.
	#
	# An alternate and simpler-looking approach would
	# be to cat all of the files into a single crd.yaml.
	find $(CRD_DIR) -type f -name '*.yaml' | \
		while read filename ; do cat $$filename > \
		$(STACK_PACKAGE_REGISTRY)/resources/$$( basename $$(echo $$filename | sed s/.yaml$$/.crd.yaml/)) \
		; done

	cp -r $(STACK_PACKAGE_REGISTRY_SOURCE)/* $(STACK_PACKAGE_REGISTRY)
.PHONY: bundle

# A local docker registry can be used to publish and consume images locally, which
# is convenient during development, as it simulates the whole lifecycle of the
# Stack, end-to-end.
docker-local-registry:
	[ $$( docker ps --filter name=registry --filter status=running --last 1 --quiet | wc -l ) -eq 1 ] || \
		docker run -d -p 5000:5000 --restart=always --name registry registry:2
.PHONY: docker-local-registry

# Tagging the image with the address of the local registry is a necessary step of
# publishing the image to the local registry.
docker-local-tag: docker-local-registry
	docker tag ${STACK_IMG} localhost:5000/${STACK_IMG}
.PHONY: docker-local-tag

# When we are developing locally, this target will publish our container image
# to the local registry.
docker-local-push: docker-local-tag docker-local-registry
	docker push localhost:5000/${STACK_IMG}
.PHONY: docker-local-push

# Sooo ideally this wouldn't be a single line, but the idea here is that when we're
# developing locally, we want to use our locally-published container image for the
# Stack's controller container. The container image is specified in the install
# yaml for the Stack. This means we need two versions of the install:
#
# * One for "regular", production publishing
# * One for development
#
# The way this has been implemented here is by having the base install yaml be the
# production install, with a yaml patch on the side for use during development.
# *This* make recipe creates the development install yaml, which requires doing
# some patching of the original install and putting it back in the stack package
# directory. It's done here as a post-processing step after the stack package
# has been generated, which is why the output to a copy and then rename step is
# needed. This is not the only way to implement this functionality.
#
# The implementation here is general, in the sense that any other yamls in the
# overrides directory will be patched into their corresponding files in the
# stack package. It assumes that all of the yamls are only one level deep.
stack-local-build: bundle
	find $(LOCAL_OVERRIDES_DIR) -maxdepth 1 -type f -name '*.yaml' | \
		while read filename ; do \
			kubectl patch --dry-run --filename $(STACK_PACKAGE_REGISTRY)/$$( basename $$filename ) \
				--type strategic --output yaml --patch "$$( cat $$filename )" > $(STACK_PACKAGE_REGISTRY)/$$( basename $$filename ).new && \
			mv $(STACK_PACKAGE_REGISTRY)/$$( basename $$filename ).new $(STACK_PACKAGE_REGISTRY)/$$( basename $$filename ) \
		; done
.PHONY: stack-local-build

# Convenience for building a local bundle by running the steps in the right order and with a single command
local-build: stack-local-build docker-build docker-local-push
.PHONY: local-build

# Install a locally-built stack using the sample stack installation CR
stack-install:
	kubectl apply -f $(CONFIG_SAMPLES_DIR)/local.install.stack.yaml
.PHONY: stack-install

stack-uninstall:
	kubectl delete -f $(CONFIG_SAMPLES_DIR)/local.install.stack.yaml
.PHONY: stack-uninstall

# Build the docker image
docker-build: bundle
	docker build --file stack.Dockerfile . -t ${STACK_IMG}
.PHONY: docker-build

docker-push:
	docker push ${STACK_IMG}
.PHONY: docker-push
