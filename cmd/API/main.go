package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"news-forex/model"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func main() {
	// สร้าง Echo instance
	e := echo.New()

	e.GET("/check-news/:dateCheck/:hour", func(c echo.Context) error {

		dateCheck := c.Param("dateCheck")
		hoursToAddStr := c.Param("hour")
		intHour, err := strconv.Atoi(hoursToAddStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid hourBefore parameter",
				"err":     err.Error(),
			})
		}

		dateTime, err := time.Parse("2006-01-02_15:04", dateCheck)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "ไม่สามารถแปลง เวลาได้",
				"err":     err.Error(),
			})
		}

		startDate := dateTime
		endDate := dateTime.Add(time.Duration(intHour) * time.Hour)
		beforeDate := dateTime.Add(-time.Duration(intHour) * time.Hour)

		data, err := ioutil.ReadFile("news_forex.json")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "ไม่สามารถอ่านไฟล์ JSON ได้",
				"err":     err.Error(),
			})
		}
		var news []model.NewsEvent
		err = json.Unmarshal(data, &news)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "ไม่สามารถแปลง JSON เป็น struct ได้",
				"err":     err.Error(),
			})
		}
		var isStatus bool

		for _, detail := range news {

			if detail.Date != dateTime.Format("02-01-2006") {
				continue
			}

			detailTime, err := time.Parse("02-01-2006 15:04", detail.Date+" "+detail.Time)

			if err != nil {
				fmt.Println("Error parsing detail time:", err)
			}

			if detailTime.Unix() >= startDate.Unix() && detailTime.Unix() <= endDate.Unix() || detailTime.Unix() <= beforeDate.Unix() {

				return c.JSON(http.StatusOK, true)
			}

			if detail.Time == model.AllDay {

				isStatus = true
			}

		}

		return c.JSON(http.StatusOK, isStatus)

	})

	e.Logger.Fatal(e.Start(":8080"))
}
