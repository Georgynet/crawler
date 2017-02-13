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

        linksRaw := linksRegExp.FindAllStringSubmatch(bodyString, -1)
        links := make([]string, len(linksRaw))

        for i, item := range linksRaw {
            links[i] = item[1]
        }

        fmt.Println(links)
    } else {
        // TODO: insert into DB: link, status, source page
    }
}
