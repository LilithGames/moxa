.PHONY: proto
proto:
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./master_shard/state.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. --dragonboat_out=paths=source_relative:. ./master_shard/api.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/member.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/event.proto
	@protoc -I. -Iproto --go_out=paths=source_relative:. ./cluster/snapshot.proto
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
	@ko build -B github.com/LilithGames/moxa/cmd/moxa

.PHONY: build
build: proto
	@go build -o bin/ github.com/LilithGames/moxa/cmd/...

.PHONY: run
run:
	@docker run -it --rm --entrypoint=bash $$(ko build -B github.com/LilithGames/moxa/cmd/moxa)

.PHONY: install
install:
	@kubectl kustomize deploy | ko resolve -B -f - | kubectl apply -f -

.PHONY: clean
clean:
	@kubectl delete -k deploy

.PHONY: deploy
deploy: install
