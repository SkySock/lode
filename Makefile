.PHONY: build-all

SERVICES = user-service
LIBRARIES = utils shared-dto


test-libs:
	@for lib in $(LIBRARIES); do \
		echo "Testing library: $$lib"; \
		go test -cover ./libs/$$lib/...; \
	done