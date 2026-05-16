module github.com/AatirNadim/getMe/http-proxy-go

go 1.23.1

require (
	github.com/AatirNadim/getMe/commons v0.0.0
	github.com/AatirNadim/getMe/sdks/goSdk v0.0.0
)

require (
	github.com/AatirNadim/getMe/utils v0.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

replace (
	github.com/AatirNadim/getMe/commons => ../commons
	github.com/AatirNadim/getMe/sdks/goSdk => ../sdks/goSdk
	github.com/AatirNadim/getMe/utils => ../utils
)
