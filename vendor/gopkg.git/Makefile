DOCKER_VERSION = 1.12.6
DOCKER_IMAGE   = golang
PROJECT_PATH   = gitlab.vpc-xpto-01/oss/gopkg.git

.PHONY: test run

run:
	@if [ "$(CI)" != "" ]; then\
		eval $(ARGS);\
	else \
		docker run \
			-t \
			--rm \
			-v "$(PWD)":/go/src/$(PROJECT_PATH) \
			-w /go/src/$(PROJECT_PATH) \
			-e "TERM=xterm-256color" \
			$(DOCKER_IMAGE):$(DOCKER_VERSION) \
			$(shell echo $(ARGS)); \
	fi

test:
	@$(MAKE) --no-print-directory \
		ARGS="go test -v -race ./..." \
		run