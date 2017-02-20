package main

import (
	"os"
	"sitemap/console"
)

func main() {
	app := console.InitApp()
	app.Run(os.Args)
}
