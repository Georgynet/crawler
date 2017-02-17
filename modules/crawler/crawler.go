package crawler

import (
	"net/http"
	"sitemap/common"
	"github.com/gin-gonic/gin"
	"github.com/golang-collections/collections/stack"
	"log"
	"io/ioutil"
	"regexp"
	"github.com/golang-collections/collections/set"
	"strings"
	"io"
	"net/url"
)

var LinksStack = stack.New()
var ExternalLinks = set.New()
var ResultLinks = set.New()
var StartUrl *url.URL

func Run(c *gin.Context, url string, sourceUrl string) {
	log.Println("[CRAWLER] Request URL: " + url)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		common.ErrorJSON(c, http.StatusBadRequest, "Error request: " + err.Error())
		return
	}
	defer resp.Body.Close()

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

		var newLink common.Link
		for _, item := range linksRaw {
			newLink = common.Link{
				Link: item[1],
				Source: url,
			}

			nLink, err := linkToAbs(newLink)
			if err != nil {
				continue
			}
			LinksStack.Push(nLink)
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

func linkToAbs(link common.Link) (common.Link, error) {
	parseUrl, err := url.Parse(link.Link)
	if err != nil {
		return link, err
	}

	if !parseUrl.IsAbs() {
		if len(link.Link) <= 0 {
			return link, *new(error)
		}

		if string(link.Link[0]) == "/" {
			link.Link = StartUrl.Scheme + "://" + StartUrl.Host + parseUrl.String()
		} else if string(link.Link[0]) == "#" {
			return link, *new(error)
		} else if strings.TrimSpace(parseUrl.String()) != "" {

			partsLink := strings.Split(parseUrl.String(), "/")
			counter := 0
			for _, item := range partsLink {
				if item == ".." {
					counter += 1
				}
			}

			partsLink = partsLink[counter:]

			partsSourceUrl := strings.Split(link.Source, "/")

			if len(partsSourceUrl) < (1 + counter) {
				return link, *new(error)
			}

			partsSourceUrl = partsSourceUrl[:len(partsSourceUrl) - (1 + counter)]
			link.Link = strings.Join(partsSourceUrl, "/") + "/" + strings.Join(partsLink, "/")
		}
	}

	return link, nil
}
