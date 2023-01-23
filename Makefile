.PHONY: rebuild
rebuild:
	docker build -t ngoctd/ecommerce-search:latest . && \
	docker push ngoctd/ecommerce-search

.PHONY: redeploy
redeploy:
	kubectl rollout restart deployment depl-search

.PHONY: protogen
protogen:
	protoc --proto_path=proto proto/search_service.proto proto/general.proto \
	--go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative