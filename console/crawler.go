package console

import (
	"net/http"
	"log"
	"io/ioutil"
	"regexp"
	"strings"
	"io"
	"net/url"
	"errors"
)

func RunCrawler(url string, sourceUrl string) {
	log.Println("[CRAWLER] Request URL: " + url)

	bodyByte, statusCode, err := getBody(url)

	ResultLinks.Insert(Page{
		Link: url,
		Source: sourceUrl,
		Status: statusCode,
	})

	if err != nil {
		log.Println("[CRAWLER] URL don't scaned")
		return
	}

	log.Println("[CRAWLER] URL scaned")

	bodyString := string(bodyByte[:])
	linksRegExp := regexp.MustCompile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)
	linksRaw := linksRegExp.FindAllStringSubmatch(bodyString, -1)

	var newLink Link
	for _, item := range linksRaw {
		newLink = Link{
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

func getBody(url string) ([]byte, int, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		// TODO: write to log get url errors
		return nil, resp.StatusCode, errors.New("Url return error")
	}
	defer resp.Body.Close()

	contentType, err := getRespContentType(resp)
	if err != nil {
		// TODO: write to log get contentType errors
		return nil, resp.StatusCode, errors.New("Can't get content type")
	}

	if http.StatusOK == resp.StatusCode && "text/html" == contentType {
		bodyByte, _ := ioutil.ReadAll(resp.Body)
		return bodyByte, resp.StatusCode, nil
	}

	return nil, resp.StatusCode, errors.New("Can't get page content")
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

func linkToAbs(link Link) (Link, error) {
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
