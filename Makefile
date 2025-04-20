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
	cp configs/ports.env.example configs/ports.env
	cp configs/postgres.env.example configs/postgres.env
	cp configs/redis.env.example configs/redis.env
	cp configs/kafka.env.example configs/kafka.env
	cp configs/shortener.env.example configs/shortener.env
	cp configs/tg_bot.env.example configs/tg_bot.env
	cp configs/user_service.env.example configs/user_service.env
	cp configs/analytics.env.example configs/analytics.env

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
