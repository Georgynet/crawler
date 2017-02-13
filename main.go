package main

import (
    "sitemap/modules/config"
    "sitemap/modules/prepare"
    "strconv"
)

func main() {
    //Configure Gin services
    r := config.InitGin()
    prepare.Register(r)

    r.Run(":" + strconv.Itoa(config.MainConfig.Server.ServerPort))
}
