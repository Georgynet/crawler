package console

import (
	"github.com/golang-collections/collections/set"
	"github.com/urfave/cli"
	"errors"
	"net/url"
	"github.com/golang-collections/collections/stack"
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
	prepareUrl, prepareErr := parseStartUrl(c.String("url"))
	if prepareErr != nil {
		return prepareErr
	}
	StartUrl = prepareUrl

	LinksStack.Push(Link{
		Link: c.String("url"),
		Source: "",
	})

	for LinksStack.Len() > 0 {
		analyse(LinksStack.Pop().(Link))
	}

	return nil
}

func parseStartUrl(inputUrl string) (*url.URL, error) {
	if "" == inputUrl {
		return nil, errors.New("URL isn't set.")
	}

	parseUrl, parseErr := url.Parse(inputUrl)
	if parseErr != nil {
		return nil, errors.New("Incorrect URL.")
	}

	return parseUrl, nil
}

func analyse(link Link) {
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
		ResultLinks.Insert(Page{
			Link: parseUrl.String(),
			Source: link.Source,
			Type: "ext",
			Status: 0,
		})
	}

	return
}
