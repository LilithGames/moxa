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
build-image: build-linux 
	@docker-compose -f deploy/docker-compose.yaml build moxa

.PHONY: push-image
push-image: build-image
	@docker-compose -f deploy/docker-compose.yaml push moxa

.PHONY: run-image
run-image: build-image
	@docker-compose -f deploy/docker-compose.yaml run moxa

.PHONY: build-linux
build-linux: proto
	@GOOS=linux go build -o bin/ github.com/LilithGames/moxa/cmd/...

.PHONY: build
build: proto
	@go build -o bin/ github.com/LilithGames/moxa/cmd/...

.PHONY: run
run: build-image
	@docker-compose -f deploy/docker-compose.yaml run moxa ./moxa

.PHONY: install
install:
	@kubectl apply -k deploy

.PHONY: clean
clean:
	@kubectl delete -k deploy

.PHONY: deploy
deploy: push-image install
	@kubectl rollout restart statefulset.apps/moxa -n temp
	@kubectl rollout status statefulset.apps/moxa -n temp
