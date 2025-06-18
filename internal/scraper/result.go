package scraper

type ScrapedResult struct {
	URL string
	Links []Link
	Err error
}

