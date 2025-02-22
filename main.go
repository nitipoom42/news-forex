package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/labstack/echo/v4"
)

const (
	curRencyUSD = "USD"

	impactHigh    = "icon--ff-impact-red"
	impactMedium  = "icon--ff-impact-yel"
	impactLow     = "icon--ff-impact-ora"
	impactUnknown = "icon--ff-impact-gra"

	resImpactHigh    = "High"
	resImpactMed     = "Medium"
	resImpactLow     = "Low"
	resImpactUnknown = "Unknown"
)

type NewsEvent struct {
	Date string `json:"date"`
	Time string `json:"time"`
}

func getNewsForex(dateTime string) NewsEvent {
	log.Println("Start")
	c := colly.NewCollector(
		colly.AllowedDomains("www.forexfactory.com"),
	)
	var lastDate string
	var lastTime string
	var resDayNews NewsEvent

	c.OnHTML("table.calendar__table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			date := row.ChildText("td.calendar__date")
			curRency := row.ChildText("td.calendar__currency")
			newTime := row.ChildText("td.calendar__time")

			if date != "" {
				lastDate = date
			} else {
				date = lastDate
			}

			if newTime != "" {
				lastTime = newTime
			} else {
				newTime = lastTime
			}
			impactClass := row.ChildAttr("td.calendar__impact span", "class")

			if curRency == curRencyUSD {
				dateForWeb, _ := convertToUnix(date, "web")
				dateForCheck, _ := convertToUnix(dateTime, "check")

				if dateForWeb == dateForCheck {
					resCheck := checkImpact(curRency, date, dateTime, impactClass)
					if resCheck == true {
						NewsEvent := NewsEvent{
							Date: date,
							Time: newTime,
						}
						resDayNews = NewsEvent
					}
				}

			}
		})
	})

	err := c.Visit("https://www.forexfactory.com/calendar?week=" + dateTime)
	if err != nil {
		log.Println("Error visiting Forex Factory:", err)
	}

	return resDayNews

}

func checkImpact(curRency, date, dateTime, impactClass string) bool {
	if curRency == curRencyUSD {
		dateForWeb, _ := convertToUnix(date, "web")
		dateForCheck, _ := convertToUnix(dateTime, "check")

		if dateForWeb == dateForCheck {
			impact := mapImpact(impactClass)
			if impact == resImpactHigh {
				return true
			}
		}
	}
	return false
}

func formatDate(input string) (string, error) {
	t, err := time.Parse("2006-01-02 15:04", input)
	if err != nil {
		return "", err
	}

	return t.Format("Jan02.2006"), nil
}

func convertToUnix(input string, typeDate string) (int64, error) {
	var layout string
	if typeDate == "web" {
		layout = "Mon Jan 2"
	} else if typeDate == "check" {
		layout = "Jan02.2006"
	}
	t, err := time.Parse(layout, input)
	if err != nil {
		return 0, err
	}

	currentYear := time.Now().Year()
	t = time.Date(currentYear, t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)

	return t.Unix(), nil
}

func mapImpact(class string) string {
	switch {
	case strings.Contains(class, impactHigh):
		return resImpactHigh
	case strings.Contains(class, impactMedium):
		return resImpactMed
	case strings.Contains(class, impactLow):
		return resImpactLow
	default:
		return resImpactUnknown
	}
}

func handleGetNews(c echo.Context) error {

	dateTime := c.Param("dateTime")
	formattedDate, err := formatDate(dateTime)
	if err != nil {
		panic(err)
	}

	resNews := getNewsForex(formattedDate)

	return c.JSON(http.StatusOK, resNews)
}

func setupRoutes(e *echo.Echo) {
	e.GET("/news/usd/:dateTime", handleGetNews)
}

func main() {
	e := echo.New()
	setupRoutes(e)
	currentTime := time.Now()
	rangDate := currentTime.Format("Jan02.2006")

	go func() {
		newsForex := getNewsForex(rangDate)
		data, err := json.MarshalIndent(newsForex, "", "  ")
		if err != nil {
			log.Println("Error marshaling JSON:", err)
			return
		}

		err = os.WriteFile("news_forex.json", data, 0644)
		if err != nil {
			log.Println("Error writing to file:", err)
			return
		}

		log.Println("Forex news has been written to news_forex.json")
	}()

	log.Println("Server is running on :8080")
	e.Logger.Fatal(e.Start(":8080"))

}
