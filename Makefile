PWD = $(shell pwd)
PROJ_DIR = $(PWD)

gen_graphql:
	@cd $(PROJ_DIR)/graphql && go generate ./...