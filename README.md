# Amazon Reviews Crawler Go

## Summary

Simple tool to crawl Amazon Reviews metadata from product ID(s).
If running via terminal, simply pass the list of product ID(s) as argument, and pass the optioonal `-out=output.json` flag to export resulted json to `output.json` file.

If you want to crawl from different Amazon regional website, you can pass either `uk` or `us` optional region flag, if region mapping isn't found, region flag input will be used as website domain suffix.

Also, you may need to add the `-ref` flag to maintain a session to avoid empty results from amazon, e.g. `-ref="sr_1_2?keywords=kano&qid=1553793792&s=gateway&sr=8-2"`

First iteration is scripted in Python and is located at [here](https://github.com/anzellai/amazon_reviews_crawler)


## Setup

`go get`


## Example

Crawling product ID "ABCDE", and run:

`go run cmd/crawler/crawler.go -region=uk -out=output.json ABCDE`


Or, run existing pre-built binary from `/bin` directory:

`./bin/amazon-reviews-crawler-go-darwin -region=uk -out=output.json ABCDE`
