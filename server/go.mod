module github.com/AatirNadim/getMe/server

go 1.23.1

require github.com/AatirNadim/getMe/commons v0.0.0

require (
	github.com/AatirNadim/getMe/utils v0.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0
)

replace (
	github.com/AatirNadim/getMe/commons => ../commons
	github.com/AatirNadim/getMe/utils => ../utils
)
