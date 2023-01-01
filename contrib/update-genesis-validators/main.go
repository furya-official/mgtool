package main

import (
	"github.com/mage-labs/mage/app"
	"github.com/mage-labs/mgtool/contrib/update-genesis-validators/cmd"
)

func main() {
	app.SetSDKConfig()
	cmd.Execute()
}
