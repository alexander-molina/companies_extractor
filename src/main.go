package main

func main() {
	var scrapers = [...]func(){Wadline}
	for _, scraper := range scrapers {
		scraper()
	}
}
