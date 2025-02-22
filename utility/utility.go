package utility

import (
	"log"
	"news-forex/model"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

func GetNewsForex(dateTime string, wg *sync.WaitGroup, mu *sync.Mutex, resDayNews *[]model.NewsEvent) {
	defer wg.Done()
	c := colly.NewCollector(
		colly.AllowedDomains("www.forexfactory.com"),
	)
	var lastDate string
	var lastTime string

	c.OnHTML("table.calendar__table", func(e *colly.HTMLElement) {
		e.ForEach("tr", func(i int, row *colly.HTMLElement) {
			date := row.ChildText("td.calendar__date")
			curRency := row.ChildText("td.calendar__currency")
			newTime := row.ChildText("td.calendar__time")
			title := row.ChildText("td.calendar__event")

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

			if curRency == model.CurRencyUSD {
				dateForWeb, _ := convertToUnix(date, "web")
				dateForCheck, _ := convertToUnix(dateTime, "check")

				if dateForWeb == dateForCheck {
					resCheck := checkImpact(curRency, date, dateTime, impactClass)
					if resCheck == true {

						time, errTime := ConvertTo24HourFormat(newTime)
						if errTime != nil {
							time = model.AllDay
						}

						NewsEvent := model.NewsEvent{
							Date:  ConvertUnixToDate(dateForCheck),
							Time:  time,
							Title: title,
						}
						mu.Lock()
						*resDayNews = append(*resDayNews, NewsEvent)
						mu.Unlock()
					}
				}

			}
		})
	})

	err := c.Visit("https://www.forexfactory.com/calendar?day=" + dateTime)
	if err != nil {
		log.Println("Error visiting Forex Factory:", err)
	}

}

func checkImpact(curRency, date, dateTime, impactClass string) bool {
	if curRency == model.CurRencyUSD {
		dateForWeb, _ := convertToUnix(date, "web")
		dateForCheck, _ := convertToUnix(dateTime, "check")
		if dateForWeb == dateForCheck {
			impact := mapImpact(impactClass)
			if impact == model.ResImpactHigh {
				return true
			}
		}
	}
	return false
}

func ConvertUnixToDate(unixTimestamp int64) string {
	t := time.Unix(unixTimestamp, 0).Local()
	return t.Format("02-01-2006")
}

func ConvertTo24HourFormat(timeString string) (string, error) {
	timeString = strings.ToUpper(timeString)

	t, err := time.Parse("3:04PM", timeString)
	if err != nil {
		return "", err
	}

	return t.Format("15:04"), nil
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
	case strings.Contains(class, model.ImpactHigh):
		return model.ResImpactHigh
	case strings.Contains(class, model.ImpactMedium):
		return model.ResImpactMed
	case strings.Contains(class, model.ImpactLow):
		return model.ResImpactLow
	default:
		return model.ResImpactUnknown
	}
}
