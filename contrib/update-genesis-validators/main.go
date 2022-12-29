package main

import (
	"github.com/furya-official/mage/app"
	"github.com/furya-official/mgtool/contrib/update-genesis-validators/cmd"
)

func main() {
	app.SetSDKConfig()
	cmd.Execute()
}
