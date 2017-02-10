package prepare

import (
    "github.com/gin-gonic/gin"
    "net/url"
    "net/http"
    "sitemap/modules/common"
)

func Parse(c *gin.Context) {
    rawUrl := c.PostForm("url")
    if rawUrl == "" {
        common.ErrorJSON(c, http.StatusBadRequest, "URL isn't set")
        return
    }

    parseUrl, parseErr := url.Parse(rawUrl)
    if parseErr != nil {
        common.ErrorJSON(c, http.StatusBadRequest, parseErr.Error())
        return
    }

    internalHost := c.PostForm("internalHost")
    if internalHost == "" {
        common.ErrorJSON(c, http.StatusBadRequest, "InternalHost isn't set")
        return
    }

    if parseUrl.Host == internalHost {
        // TODO: start crawler
        c.JSON(http.StatusOK, gin.H{
            "status":  "internal",
            "message": "parseUrl.Host == internalHost",
        })
        return
    } else {
        // TODO: insert into DB
        c.JSON(http.StatusOK, gin.H{
            "status":  "external",
            "message": "parseUrl.Host != internalHost",
        })
        return
    }
}