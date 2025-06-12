.PHONY: build-all

SERVICES = user-service
LIBRARIES = utils shared-dto


test-libs:
	@for lib in $(LIBRARIES); do \
		echo "Testing library: $$lib"; \
		go test -cover ./libs/$$lib/...; \
	done

gen-swagger:
	@for service in $(SERVICES); do \
		echo "Generate swagger: $$service"; \
		swag init -g ./services/$$service/cmd/main.go -o ./services/$$service/docs --parseDependency --parseInternal; \
	done