package console

import (
	"os"
	"log"
	"fmt"
	"github.com/golang-collections/collections/set"
)

// Page of site
type Page struct {
	Link   string
	Source string
	Type   string
	Status int
}

// Visited link
type Link struct {
	Link   string
	Source string
}

// Save visited links to file
func SaveVisitedLinks(links *set.Set, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	links.Do(func(item interface{}) {
		fmt.Fprintln(file, item.(string))
	})
}

// Save Page-struct to file
func SaveResultLinks(pages *set.Set, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintln(file, "link;source;type;status")
	pages.Do(func(item interface{}) {
		fmt.Fprintln(file, item.(Page).Link + ";" + item.(Page).Source + ";" + item.(Page).Type + ";" + fmt.Sprint(item.(Page).Status))
	})
}
