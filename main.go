package main

import (
	"os"
	"sitemap/modules/console"
)

func main() {
	app := console.InitApp()
	app.Run(os.Args)
}
