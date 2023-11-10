# Makefile

PROJECT_DIR=$GOPATH/src/github.com/bazel/bazel_grpc_pbuf_go
PROTO_DIR=$(PROJECT_DIR)/proto
BZCMD=bazel7

.PHONY:build_all build_client build_proto build_server clean generate_repos go_mod_tidy list run_client run_server test update_repos

build_all: build_proto build_server run_client
	echo "DONE"

build_client:
	@$(BZCMD) build --sandbox_debug --verbose_failures //cmd/client:client

build_proto:
	@$(BZCMD) build //proto/helloworld/v1:helloworld_v1_go

build_server:
	@$(BZCMD) build --sandbox_debug --verbose_failures //cmd/server:server

clean:
	@$(BZCMD) clean --expunge --async

generate_repos:
	@$(BZCMD) run //:gazelle

go_mod_tidy:
	@$(BZCMD) run @io_bazel_rules_go//go -- mod tidy

list:
	@$(BZCMD) query //...

run_client:
	@$(BZCMD) run --sandbox_debug --verbose_failures //cmd/client:client

run_server:
	@$(BZCMD) run --sandbox_debug --verbose_failures //cmd/server:server

test:
	@$(BZCMD) test //...

update_repos:
	@$(BZCMD) run //:gazelle-update-repos -- -from_file=go.mod -to_macro=deps.bzl%go_dependencies -prune
