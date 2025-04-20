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

copy_env:
	cp ports.env.example ports.env
	cp postgres.env.example postgres.env
	cp redis.env.example redis.env
	cp shortener.env.example shortener.env
	cp tg_bot.env.example tg_bot.env
	cp token_service.env.example token_service.env

.PHONY: build-all
build-all:
	@docker compose build token_service &
	@docker compose build shortener &
	@docker compose build tg_bot &
	wait

env_for_build:
	touch build/docker/.env
	grep -E '^(TOKENS_PORT_GRPC|SHORTENER_PORT_GRPC|POSTGRES_PORT|REDIS_PORT)=' \
	configs/ports.env > build/docker/.env
