package main

import (
	"github.com/jucardi/infuse/cmd/infuse/cli"
	_ "github.com/jucardi/infuse/templates/gotmpl"
)

func main() {
	cli.Execute()
}
