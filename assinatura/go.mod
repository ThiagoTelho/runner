module github.com/thiagotelho/runner/assinatura

go 1.21.3

require github.com/spf13/cobra v1.10.2

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/thiagotelho/runner/jdk v0.0.0
)

replace github.com/thiagotelho/runner/jdk => ../jdk
