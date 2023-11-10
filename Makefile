# Makefile

PROJECT_DIR=$GOPATH/src/github.com/bazel/easyproto2
PROTO_DIR=$(PROJECT_DIR)/proto
BZCMD=bazel

.PHONY:baz_build_all baz_clean baz_go_command baz_list baz_run_client baz_tidy

baz_tidy: baz_clean
	@rm proto/abitofhelp/api/blob/v5/BUILD.bazel \
        proto/abitofhelp/enum/v5/BUILD.bazel

#baz_build_enum_lib:
#	@$(BZCMD) build //src/go/gen/abitofhelp/enum/v5:enum --sandbox_debug --verbose_failures

baz_go_mod_tidy:
	@$(BZCMD) run go -- mod tidy

baz_build_proto:
	@$(BZCMD) build //proto/helloworld/v1:helloworld_v1_go

baz_build_client:
	@$(BZCMD) build --sandbox_debug --verbose_failures //cmd/client:client

baz_build_server:
	@$(BZCMD) build --sandbox_debug --verbose_failures //cmd/server:server

baz_build_all:
	@$(BZCMD) build //... #$(BZCMD) build //src/client:client

baz_run_client:
	@$(BZCMD) run --sandbox_debug --verbose_failures //cmd/client:client

baz_run_server:
	@$(BZCMD) run --sandbox_debug --verbose_failures //cmd/server:server


baz_clean:
	@$(BZCMD) clean --expunge --async

baz_list:
	@$(BZCMD) query //...

baz_go_command:
	@$(BZCMD) run @io_$(BZCMD)_rules_go//go -- build ./...

baz_test:
	@$(BZCMD) test //...

1.gazelle_generate_repos:
	@$(BZCMD) run //:gazelle

2.gazelle_update_repos:
	@$(BZCMD) run //:gazelle-update-repos -- -from_file=go.mod -to_macro=deps.bzl%go_dependencies -prune
