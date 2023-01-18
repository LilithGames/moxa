define KO
MVERSION=$$(go run github.com/LilithGames/moxa/tools/pkghash -short -pkg github.com/LilithGames/moxa/master_shard) \
SVERSION=$$(go run github.com/LilithGames/moxa/tools/pkghash -short -pkg github.com/LilithGames/moxa/sub_shard) \
ko
endef

.PHONY: proto
proto:
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./master_shard/state.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. --dragonboat_out=paths=source_relative:. ./master_shard/api.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/member.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/event.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/snapshot.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/config.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./service/config.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./service/api.proto
	@protoc -I. -Iproto --grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true ./service/api.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. --dragonboat_out=paths=source_relative:. ./sub_shard/api.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./sub_service/api.proto
	@protoc -I. -Iproto --grpc-gateway_out=. --grpc-gateway_opt=logtostderr=true --grpc-gateway_opt=paths=source_relative --grpc-gateway_opt=generate_unbound_methods=true ./sub_service/api.proto

.PHONY: clean-proto
clean-proto:
	@rm -f master_shard/*.pb*.go
	@rm -f sub_shard/*.pb*.go
	@rm -f cluster/*.pb*.go
	@rm -f service/*.pb*.go

.PHONY: build-image
build-image: proto
	@$(KO) build -B github.com/LilithGames/moxa/cmd/moxa


.PHONY: build-linux
build-linux: proto
	@GOOS=linux go build -o bin/ github.com/LilithGames/moxa/cmd/...

.PHONY: build-base-mage
build-base-mage: build-linux
	@docker-compose -f deploy/docker-compose.yaml build moxactl

.PHONY: build
build: proto
	@go build -o bin/ github.com/LilithGames/moxa/cmd/...

.PHONY: run
run: build-base-mage
	@docker run -it --rm --entrypoint=bash $$($(KO) build -B github.com/LilithGames/moxa/cmd/moxa)

.PHONY: run-ctl
run-ctl: build
	@bin/moxactl.exe $@

.PHONY: install
install: build-base-mage
	@kubectl kustomize deploy | $(KO) resolve -B -f - | kubectl apply -f -

.PHONY: clean
clean:
	@kubectl delete -k deploy

.PHONY: deploy
deploy: install
	@kubectl rollout status statefulset.apps/moxa -n temp
