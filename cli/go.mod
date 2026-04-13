module github.com/AatirNadim/getMe/cli

go 1.23.1

require (
	github.com/AatirNadim/getMe/utils v0.0.0-20260412185038-90ed5e87cf59
	github.com/spf13/cobra v1.10.1
	github.com/spf13/pflag v1.0.9
)

require (
	github.com/AatirNadim/getMe/commons v0.0.0-00010101000000-000000000000
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
)

replace github.com/AatirNadim/getMe/commons => ../commons
