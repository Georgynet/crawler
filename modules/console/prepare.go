package console

import (
	"github.com/golang-collections/collections/set"
	"github.com/urfave/cli"
	"errors"
	"net/url"
	"github.com/golang-collections/collections/stack"
	"sitemap/modules/common"
	"sitemap/modules/crawler"
)

// Set of visited links
var VisitedLinks = set.New()
// Start analyse URL
var StartUrl *url.URL
// Links for analyse
var LinksStack = stack.New()
// Result links
var ResultLinks = set.New()

func Parse(c *cli.Context) error {
	prepareUrl, prepareErr := prepare(c.String("url"))
	if prepareErr != nil {
		return prepareErr
	}
	StartUrl = prepareUrl

	LinksStack.Push(common.Link{
		Link: c.String("url"),
		Source: "",
	})

	for LinksStack.Len() > 0 {
		analyse(crawler.LinksStack.Pop().(common.Link))
	}

	return nil
}

func prepare(inputUrl string) (*url.URL, error) {
	if "" == inputUrl {
		return nil, errors.New("URL isn't set.")
	}

	parseUrl, parseErr := url.Parse(inputUrl)
	if parseErr != nil {
		return nil, errors.New("Incorrect URL.")
	}

	return parseUrl, nil
}

func analyse(link common.Link) {
	if VisitedLinks.Has(link.Link) {
		return
	}
	VisitedLinks.Insert(link.Link);

	parseUrl, parseErr := url.Parse(link.Link)
	if parseErr != nil {
		// TODO: write to log url parse errors
		return
	}

	if parseUrl.Host == StartUrl.Host {
		// run crawler
	} else {
		ResultLinks.Insert(common.Page{
			Link: parseUrl.String(),
			Source: link.Source,
			Type: "ext",
			Status: 0,
		})
	}

	return
}
