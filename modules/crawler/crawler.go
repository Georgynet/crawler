package crawler

import (
    "net/http"
    "sitemap/modules/common"
    "github.com/gin-gonic/gin"
)

func Crawler(c *gin.Context, url string) {
    resp, err := http.Get(url)
    if err != nil {
        common.ErrorJSON(c, http.StatusBadRequest, "Error request")
        return
    }

    if http.StatusOK == resp.StatusCode {
        // TODO:
    }
}
