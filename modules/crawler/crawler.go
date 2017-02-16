package crawler

import (
    "net/http"
    "sitemap/modules/common"
    "github.com/gin-gonic/gin"
    "github.com/golang-collections/collections/stack"
    "log"
    "io/ioutil"
    "regexp"
    "github.com/golang-collections/collections/set"
    "strings"
    "io"
)

var LinksStack = stack.New()
var ExternalLinks = set.New()
var ResultLinks = set.New()

func Run(c *gin.Context, url string, sourceUrl string) {
    log.Println("[CRAWLER] Request URL: " + url)
    resp, err := http.Get(url)
    defer resp.Body.Close()

    if err != nil {
        common.ErrorJSON(c, http.StatusBadRequest, "Error request: " + err.Error())
        return
    }

    contentType, err := getRespContentType(resp)
    if err != nil {
        common.ErrorJSON(c, http.StatusBadRequest, "Error get type content: " + err.Error())
        return
    }

    if http.StatusOK == resp.StatusCode && "text/html" == contentType {
        log.Println("[CRAWLER] Request URL: " + url + " scaned")
        bodyByte, _ := ioutil.ReadAll(resp.Body)
        bodyString := string(bodyByte[:])

        linksRegExp := regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)

        linksRaw := linksRegExp.FindAllStringSubmatch(bodyString, -1)

        for _, item := range linksRaw {
            LinksStack.Push(common.Link{
                Link: item[1],
                Source: url,
            })
        }
    }

    ResultLinks.Insert(common.Page{
        Link: url,
        Source: sourceUrl,
        Status: resp.StatusCode,
    })
}

func getRespContentType(resp *http.Response) (string, error) {
    buffer := make([]byte, 512)
    n, err := resp.Body.Read(buffer)
    if err != nil && err != io.EOF {
        return "", err
    }
    contentType := http.DetectContentType(buffer[:n])

    cType := strings.Split(contentType, ";")

    return cType[0], nil
}
