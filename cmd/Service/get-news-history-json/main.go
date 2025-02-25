package main

import (
	"encoding/json"
	"log"
	"news-forex/model"
	"news-forex/utility"
	"os"
	"sync"
	"time"
)

func main() {
	// currentTime := time.Now()
	// rangDate := currentTime.Format("Jan02.2006")

	listDate := []string{}
	startDate := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 2, 22, 0, 0, 0, 0, time.UTC)

	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		listDate = append(listDate, d.Format("Jan02.2006"))
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var newsForex []model.NewsEvent

	for _, date := range listDate {
		wg.Add(1)
		go utility.GetNewsForex(date, &wg, &mu, &newsForex)
	}

	wg.Wait()

	data, err := json.MarshalIndent(newsForex, "", "  ")
	if err != nil {
		log.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("news_forex_history.json", data, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

}
