module xlink/token_service

go 1.24.1

require (
	github.com/go-redis/redis/v7 v7.4.1
	github.com/google/uuid v1.6.0
	github.com/ilyakaznacheev/cleanenv v1.5.0
	go.uber.org/zap v1.27.0
	google.golang.org/grpc v1.71.1
	google.golang.org/protobuf v1.36.6
	xlink/common v0.0.0-00010101000000-000000000000
)

replace xlink/common => ../common

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/chempik1234/common-chempik-pkg-golang v1.0.1 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250414145226-207652e42e2e // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
