module github.com/AatirNadim/getMe/sdks/goSdk

go 1.23.1

require (
	github.com/AatirNadim/getMe/commons v0.0.0
	github.com/joho/godotenv v1.5.1
)

require github.com/AatirNadim/getMe/utils v0.0.0 // indirect

replace (
	github.com/AatirNadim/getMe/commons => ../../commons
	github.com/AatirNadim/getMe/utils => ../../utils
)
