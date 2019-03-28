package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/gocolly/colly"
)

var (
	// VERSION of this tool
	VERSION string
	// COMMIT SHA of this tool
	COMMIT string
	// BRANCH for git branch
	BRANCH string

	regionMaps = map[string]string{
		"uk": "co.uk",
		"us": "com",
	}
)

type review struct {
	ProductID string `json:"product_id"`
	Profile   string `json:"profile"`
	Star      string `json:"star"`
	Title     string `json:"title"`
	Comment   string `json:"comment"`
}

func crawl(productID, region, ref string) []*review {
	if productID == "" {
		log.Println("Amazon Reviews Product ID required")
		os.Exit(1)
	}
	if regionSuffix, ok := regionMaps[region]; ok {
		region = regionSuffix
	}
	reviews := make([]*review, 0)

	// Instantiate default collector
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36"),
		colly.Async(true),
	)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 20})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cache-Control", "no-cache")
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach(".review", func(_ int, el *colly.HTMLElement) {
			r := &review{}
			r.ProductID = productID
			r.Profile = el.ChildText(".a-profile-name")
			r.Star = el.ChildText(".review-rating span")
			r.Title = el.ChildText("a.review-title span")
			r.Comment = el.ChildText("span.review-text")
			reviews = append(reviews, r)
		})
		if link, found := e.DOM.Find(".a-last a[href]").Attr("href"); found {
			log.Println("crawling next page of reviews: ", link)
			e.Request.Visit(link)
		}
	})

	startPage := fmt.Sprintf("https://www.amazon.%s/reviews/%s", region, productID)
	if ref != "" {
		startPage = startPage + "/ref=" + ref
	}
	log.Println("crawling first reviews page: ", startPage)
	c.Visit(startPage)
	c.Wait()
	return reviews
}

func main() {
	var (
		region string
		ref    string
		out    string
	)
	flag.StringVar(&region, "region", "uk", "region or amazon domain suffix, eg. -region=uk")
	flag.StringVar(&ref, "ref", "", "reference suffix querystring from amazon reviews page ending with /ref=...")
	flag.StringVar(&out, "out", "", "JSON output to file, eg. -out=output.json")
	flag.Parse()
	log.Printf("amazon-reviews-crawler-go tool -- Version: %s, Commit: %s, Branch: %s\n", VERSION, COMMIT, BRANCH)
	log.Printf("crawling on region: '%s' with product id(s): %v, output: %s\n", region, flag.Args(), out)

	var productReviews = make([]*review, 0)
	for _, productID := range flag.Args() {
		productReviews = append(productReviews, crawl(productID, region, ref)...)
	}
	if out == "" {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(productReviews)
	} else {
		reviewsJSON, err := json.MarshalIndent(productReviews, "", "  ")
		if err != nil {
			log.Println("reviews marhsalling error: ", err.Error())
			os.Exit(1)
		}
		err = ioutil.WriteFile(out, reviewsJSON, 0644)
		if err != nil {
			log.Println("output write file error: ", err.Error())
			os.Exit(1)
		}
	}
}
