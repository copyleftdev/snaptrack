package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/copyleftdev/snaptrack/pkg/crawl"
	"github.com/copyleftdev/snaptrack/pkg/store"
	"github.com/copyleftdev/snaptrack/ui"
)

func main() {

	db, err := store.InitDB("snapshots.db")
	if err != nil {
		log.Fatal("Failed to init DB:", err)
	}
	defer db.Close()

	if len(os.Args) < 2 {

		if err := ui.StartProgram(db); err != nil {
			log.Fatal("TUI error:", err)
		}
		return
	}

	switch os.Args[1] {
	case "crawl":

		if len(os.Args) < 3 {
			fmt.Println("Usage: snapstack crawl <url> [--max-depth=N]")
			return
		}
		startURL := os.Args[2]
		config := crawl.CrawlerConfig{MaxDepth: 3, Concurrency: 5}

		for i := 3; i < len(os.Args); i++ {
			if os.Args[i] == "--max-depth" && i+1 < len(os.Args) {
				md, err := strconv.Atoi(os.Args[i+1])
				if err == nil {
					config.MaxDepth = md
				}
			}
		}

		if err := crawl.CrawlDomain(startURL, db, config); err != nil {
			fmt.Println("Crawl error:", err)
		}

	case "check":

		fmt.Println("Check not yet implemented.")

	case "diff":

		fmt.Println("Diff command not yet implemented.")

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
	}
}
