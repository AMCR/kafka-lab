
SWAGGER_VERSION=v0.31.0
CUR_DIR=$(shell pwd)

.PHONY: swagger-generate-%
swagger-generate-%:
	rm -rf ./publisher/internal/swagger-clients/archive
	mkdir -p ./publisher/internal/swagger-clients/archive
	docker run --rm --network none \
		--user "$(shell id -u):$(shell id -g)" \
		-e GOFLAGS=-mod=vendor -v "${CUR_DIR}:/app" -w /app quay.io/goswagger/swagger:$(SWAGGER_VERSION) generate client -f $*.yaml  -t ./publisher/internal/swagger-clients/archive

.PHONY: swagger-generate
swagger-generate: swagger-generate-archive-product