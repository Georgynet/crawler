package main

import (
    "sitemap/modules/config"
    "sitemap/modules/prepare"
)

func main() {
    //Configure Gin services
    r := config.InitGin()
    prepare.Register(r)

    r.Run(":3000")
}
