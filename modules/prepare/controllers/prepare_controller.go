package prepare

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"net/http"
	"sitemap/common"
	"sitemap/modules/crawler"
	"github.com/golang-collections/collections/set"
)

// Set of visited links
var VisitedLinks = set.New()

// Prepare link and run crawler
func Parse(c *gin.Context) {
	rawUrl := c.PostForm("url")
	if rawUrl == "" {
		common.ErrorJSON(c, http.StatusBadRequest, "URL isn't set")
		return
	}

	var parseErr error
	crawler.StartUrl, parseErr = url.Parse(rawUrl)
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

		if internalHost == "" {
			common.ErrorJSON(c, http.StatusBadRequest, "InternalHost isn't set")
			continue
		}

		if parseUrl.Host == internalHost {
			crawler.Run(c, parseUrl.String(), link.Source)
		} else {
			crawler.ExternalLinks.Insert(common.Page{
				Link: parseUrl.String(),
				Source: link.Source,
				Status: 0,
			})
			continue
		}
	}

	common.SaveVisitedLinks(VisitedLinks)
	common.SaveResult(crawler.ResultLinks, "pages.csv")
	common.SaveResult(crawler.ExternalLinks, "external.csv")
}
