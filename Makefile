NAMES= errors core tracing logging utils sfqb apperror metrics minio postgresql redis nats queryify


.PHONY: tags
tags: SHELL:=/bin/bash
tags:
	@for name in $(NAMES); do \
  		version=$$(cat $(CURDIR)/$$name/version) && echo "work with $$name and $$version" && \
  		tag=$$name/v$$version && echo "tag: $$tag" && \
  		if [[ ! $$(git tag -l "$$tag") ]]; then \
  		  	git tag -a "$$tag" -m "" && \
  		  	git push origin "$$tag" -o ci.skip; \
		fi; \
	done

.PHONY: init
init:
	@for name in $(NAMES); do \
  	  echo $$name && \
	  cd $$name && \
	  go mod init github.com/libraries/$$name || true && go mod tidy || true && go get -u ./... && cd ..; \
  	done
