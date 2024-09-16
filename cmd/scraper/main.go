package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joaopugsley/hltv-scraper/internal/scraper"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-site-isolation-trials", true),
		chromedp.Flag("headless", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, options...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx)
	defer cancel()

	matchesData, err := scraper.ScrapeMatches(ctx)
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.Marshal(matchesData)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("matches.json", jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Matches saved to matches.json")
}
