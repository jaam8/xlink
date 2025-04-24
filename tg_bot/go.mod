module xlink/tg_bot

go 1.24.1

require (
	github.com/ilyakaznacheev/cleanenv v1.5.0
	github.com/mymmrac/telego v1.0.2
	xlink/common v0.0.0-00010101000000-000000000000
)

replace (
	xlink/common => ../common
)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/bytedance/sonic v1.13.2 // indirect
	github.com/bytedance/sonic/loader v0.2.4 // indirect
	github.com/cloudwego/base64x v0.1.5 // indirect
	github.com/grbit/go-json v0.11.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.61.0 // indirect
	github.com/valyala/fastjson v1.6.4 // indirect
	golang.org/x/arch v0.16.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)
