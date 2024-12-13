package main

import (
	"errors"
	"fmt"
	"time"
)

type Day struct {
	date string
	fajr string
	fajrSlutt string
	duhr string
	asr string
	maghrib string
	isha string
}

type Calendar struct {
	days []Day
}

func (c *Calendar) findDay(date string) (*Day, error) {
	const layout string = "02-01-2006"
	for _, d := range c.days {
		if date == d.date {
			return &d, nil
		}
	}
	return nil, errors.New("date not found")
}

func Midnight(s string, c *Calendar) string {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return "time zone error"
	}

	const layout string = "02-01-2006 15:04"

	d, err := c.findDay(s)
	if err != nil {
		fmt.Println("Error finding day", err)
		return "err"
	}

	mdate, err := time.Parse("02-01-2006", d.date)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return "date parsing error"
	}

	maghribTime := fmt.Sprintf("%s %s", d.date, d.maghrib)
	start, err := time.ParseInLocation(layout, maghribTime, location)
	if err != nil {
		fmt.Println("Error parsing Maghrib time:", err)
		return "maghrib parsing error"
	}

	fdate := mdate.Add(24 * time.Hour).Format("02-01-2006")
	fajrDay, err := c.findDay(fdate)
	var fajrtime string

	if err == nil {
		fajrtime = fmt.Sprintf("%s %s", fdate, fajrDay.fajr)
	} else {
		fajrtime = fmt.Sprintf("%s %s", fdate, d.fajr)
	}

	end, err := time.ParseInLocation(layout, fajrtime, location)
	if err != nil {
		fmt.Println("Error parsing Fajr time:", err)
		return "fajr parsing error"
	}
	
	elapsed := end.Sub(start)
	midnight := start.Add(elapsed / 2)

	formattedMidnight := midnight.Format("15:04")

	return formattedMidnight
}
