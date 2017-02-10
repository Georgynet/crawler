package crawler

import (
    "net/http"
    "sitemap/modules/common"
    "github.com/gin-gonic/gin"
    "log"
    "fmt"
    "io/ioutil"
    "regexp"
)

func Run(c *gin.Context, url string) {
    log.Println("[CRAWLER] Request URL: " + url)
    resp, err := http.Get(url)
    if err != nil {
        common.ErrorJSON(c, http.StatusBadRequest, "Error request: " + err.Error())
        return
    }

    if http.StatusOK == resp.StatusCode {
        defer resp.Body.Close()
        bodyByte, _ := ioutil.ReadAll(resp.Body)
        bodyString := string(bodyByte[:])

        linksRegExp := regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)

        fmt.Println(linksRegExp.FindAllStringSubmatch(bodyString, -1))
    } else {
        // TODO: insert into DB: link, status, source page
    }
}
