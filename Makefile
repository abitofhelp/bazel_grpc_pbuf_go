# Makefile

PROJECT_DIR=$GOPATH/src/github.com/bazel_grpc_pbuf_go
PROTO_DIR=$(PROJECT_DIR)/proto
BZLCMD=bazel
BAZEL_BUILD_OPTS:=--verbose_failures --sandbox_debug

.PHONY:build_all build_client build_proto build_server clean gazelle_update_repos generate_repos go_mod_tidy list run_client run_server test update_repos

build_all:
	$(BZLCMD) build $(BAZEL_BUILD_OPTS) //...

build_client:
	$(BZLCMD) build $(BAZEL_BUILD_OPTS) //cmd/client:client

build_proto:
	$(BZLCMD) build $(BAZEL_BUILD_OPTS) //proto/helloworld/v1:helloworld_v1_go

build_server:
	$(BZLCMD) build $(BAZEL_BUILD_OPTS) //cmd/server:server

clean:
	$(BZLCMD) clean --expunge --async

generate_repos:
# This will generate new BUILD.bazel files for your project. You can run the same command in the future to update existing BUILD.bazel files to include new source files or options.
	$(BZLCMD) run $(BAZEL_BUILD_OPTS) //:gazelle

go_mod_tidy:
	# Assumes GO111MODULE=on
	#$(BZLCMD) run @io_bazel_rules_go//go -- mod tidy

go_mod_vendor:
	rm -rf vendor
	# Assumes GO111MODULE=on
	$(BZLCMD) run @io_bazel_rules_go//go -- mod vendor -v

go_mod_verify:
	# Assumes GO111MODULE=on
	$(BZLCMD) run @io_bazel_rules_go//go -- mod verify

list:
	$(BZLCMD) query //...

run_client:
	$(BZLCMD) run $(BAZEL_BUILD_OPTS) //cmd/client:client

run_server:
	$(BZLCMD) run $(BAZEL_BUILD_OPTS) //cmd/server:server

test:
	$(BZLCMD) test --test_output=all //...

gazelle_update_repos:
	# Import repositories from go.mod and update Bazel's macro and rules.
	$(BZLCMD) run $(BAZEL_BUILD_OPTS) //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories

update_repos: go_mod_verify go_mod_tidy gazelle_update_repos go_mod_vendor


