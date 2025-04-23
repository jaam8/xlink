generate_user_service:
	protoc -I ./user_service \
		--go_out=. \
		--go-grpc_out=. \
		./user_service/api/user_service.proto

generate_shortener:
	protoc -I ./shortener \
		--go_out=. \
		--go-grpc_out=. \
		./shortener/api/shortener.proto

generate_analytics:
	protoc -I ./analytics \
	--go_out=. \
	--go-grpc_out=. \
	./analytics/api/analytics.proto

yaml_to_env:
	cd scripts && \
	go run yaml_to_env.go

copy_env:
	cp configs/.env.example configs/.env

update_env_example:
	make yaml_to_env
	cp configs/.env configs/.env.example

.PHONY: build-all
build-all:
	@docker compose build token_service &
	@docker compose build shortener &
	@docker compose build tg_bot &
	wait

env_for_build:
	make yaml_to_env
	cp configs/.env build/docker/.env
