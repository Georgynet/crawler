package prepare

import (
    "github.com/gin-gonic/gin"
    "sitemap/modules/prepare/controllers"
)

func Register(r *gin.Engine) {
    group := r.Group("/api/v1")
    {
        group.POST("/prepare", prepare.Parse)
    }
}
