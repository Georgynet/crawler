package config

import (
    "github.com/gin-gonic/gin"
    "log"
)

func InitGin() (r *gin.Engine) {
    r = gin.Default()
    r.Use(gin.Recovery())

    log.Println("[GIN] Ready to serve")

    return r
}
