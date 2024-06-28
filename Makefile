REPO =? github.com/crossplane-contrib/function-pkl
CONTAINER_IMAGE =? ghcr.io/crossplane-contrib/function-pkl

# Branch used for Pkl Package Releases
BRANCH =? $(shell git branch --show-current)

.PHONY: pkl-resolve
pkl-resolve:
	pkl project resolve ./pkl/*/

.PHONY: pkl-release
pkl-release:
	@[ "${TAG}" ] || (error TAG is not specified && exit 1)
	$(eval RELEASE_FILES := $(shell pkl project package ./pkl/*/ | grep ${TAG}))
	@if [ -z "$(RELEASE_FILES)" ]; then \
		echo "No release files found for tag ${TAG}."; \
		exit 1; \
	fi

	gh release create ${TAG} \
	-t ${TAG} \
	-n "" \
	--target ${BRANCH} \
	--prerelease \
	--draft \
	$(RELEASE_FILES)

.PHONY: build-image
build-image:
	docker build --build-arg -t runtime .
	crossplane xpkg build -f package --embed-runtime-image=runtime -o .out/function-pkl.xpkg

.PHONY: push-image
push-image:
	crossplane xpkg push -f .out/function-pkl.xpkg ${CONTAINER_IMAGE}:${TAG}
