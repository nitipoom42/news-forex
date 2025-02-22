package main

import (
	"encoding/json"
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

	e.GET("/check-news/:dateCheck/:hourBefore/:hourAfter", func(c echo.Context) error {

		dateCheck := c.Param("dateCheck")
		hourBeforeStr := c.Param("hourBefore")
		hourBefore, err := strconv.Atoi(hourBeforeStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid hourBefore parameter",
				"err":     err.Error(),
			})
		}
		hourAfterStr := c.Param("hourAfter")
		hourAfter, err := strconv.Atoi(hourAfterStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid hourBefore parameter",
				"err":     err.Error(),
			})
		}

		dateTime, err := time.Parse("2006-01-02 15:04", dateCheck)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "ไม่สามารถแปลง เวลาได้",
				"err":     err.Error(),
			})
		}
		startTime := dateTime.Add(+time.Duration(hourBefore) * time.Hour)
		endTime := dateTime.Add(time.Duration(hourAfter) * time.Hour)

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
		var isStart bool
		var isEnd bool
		for _, detail := range news {

			newsDateTime, err := time.Parse("02-01-2006 15:04", detail.Date+" "+detail.Time)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "ไม่สามารถแปลง เวลาของข่าวได้",
					"err":     err.Error(),
				})
			}
			dateTimeNews, err := time.Parse("2006-01-02 15:04", newsDateTime.Format("2006-01-02 15:04"))
			if err != nil {
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": "ไม่สามารถแปลง เวลาของข่าวเป็น time ได้",
					"err":     err.Error(),
				})
			}

			if startTime.Hour() == dateTimeNews.Hour() && startTime.Minute() == dateTimeNews.Minute() {
				isStart = true
			}

			if endTime.Hour() == dateTimeNews.Hour() && endTime.Minute() == dateTimeNews.Minute() {
				isEnd = true
			}

		}

		return c.JSON(http.StatusOK, echo.Map{
			"startTime": isStart,
			"endTime":   isEnd,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
