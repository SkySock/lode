.PHONY: build-all

SERVICES = user-service
LIBRARIES = utils


test-libs:
	@for lib in $(LIBRARIES); do \
		echo "Testing library: $$lib"; \
		go test -cover ./libs/$$lib/...; \
	done