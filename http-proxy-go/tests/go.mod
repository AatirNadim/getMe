module github.com/AatirNadim/getMe/http-proxy-go/tests

go 1.23.1

replace (
	github.com/AatirNadim/getMe/commons => ../../commons
	github.com/AatirNadim/getMe/http-proxy-go => ../
	github.com/AatirNadim/getMe/sdks/goSdk => ../../sdks/goSdk
	github.com/AatirNadim/getMe/utils => ../../utils
)

require (
	github.com/AatirNadim/getMe/commons v0.0.0
	github.com/AatirNadim/getMe/http-proxy-go v0.0.0-00010101000000-000000000000
)

require (
	github.com/AatirNadim/getMe/sdks/goSdk v0.0.0 // indirect
	github.com/AatirNadim/getMe/utils v0.0.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)
