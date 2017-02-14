package crawler

import (
    "net/http"
    "sitemap/modules/common"
    "github.com/gin-gonic/gin"
    "github.com/golang-collections/collections/stack"
    "log"
    "io/ioutil"
    "regexp"
)

var LinksStack = stack.New()

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

        linksRaw := linksRegExp.FindAllStringSubmatch(bodyString, -1)

        for _, item := range linksRaw {
            LinksStack.Push(item[1])
        }
    } else {
        // TODO: insert into DB: link, status, source page
    }
}
