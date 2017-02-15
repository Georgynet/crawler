package prepare

import (
    "github.com/gin-gonic/gin"
    "net/url"
    "net/http"
    "sitemap/modules/common"
    "sitemap/modules/crawler"
    "github.com/golang-collections/collections/set"
    "strings"
)

var VisitedLinks = set.New()

func Parse(c *gin.Context) {
    rawUrl := c.PostForm("url")
    if rawUrl == "" {
        common.ErrorJSON(c, http.StatusBadRequest, "URL isn't set")
        return
    }

    startUrl, parseErr := url.Parse(rawUrl)
    if parseErr != nil {
        common.ErrorJSON(c, http.StatusBadRequest, "Incorrect URL")
        return
    }

    internalHost := c.PostForm("internalHost")

    crawler.LinksStack.Push(common.Link{
        Link: rawUrl,
        Source: "",
    })
    for crawler.LinksStack.Len() > 0 {

        link := crawler.LinksStack.Pop().(common.Link)

        if VisitedLinks.Has(link.Link) {
            continue
        }

        VisitedLinks.Insert(link.Link);

        parseUrl, parseErr := url.Parse(link.Link)
        if parseErr != nil {
            common.ErrorJSON(c, http.StatusBadRequest, parseErr.Error())
            continue
        }

        if !parseUrl.IsAbs() {
            crawler.LinksStack.Push(common.Link{
                Link: startUrl.Scheme + "://" + startUrl.Host + "/" + strings.TrimLeft(parseUrl.String(), "/"),
                Source: link.Source,
            })
        }

        if internalHost == "" {
            common.ErrorJSON(c, http.StatusBadRequest, "InternalHost isn't set")
            continue
        }

        if parseUrl.Host == internalHost {
            crawler.Run(c, parseUrl.String(), link.Source)
        } else {
            // TODO: insert into DB
            c.JSON(http.StatusOK, gin.H{
                "status":  "external",
                "message": "parseUrl.Host != internalHost",
            })
            continue
        }
    }

    common.SaveVisitedLinks(VisitedLinks)
    common.SaveResult(crawler.ResultLinks)
}
