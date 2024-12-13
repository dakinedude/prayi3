package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	const prefix string = "Bonnetid_Oslo_"
	const suffix string = ".csv"
	const format string = "2006-01"

	mid := time.Now().Format(format)
	fn := fmt.Sprintf("%s%s%s", prefix, mid, suffix)
	file, err := os.Open(fmt.Sprintf("/home/mats/programming/go/bonnetid/%s", fn))
	if err != nil {
		fmt.Println("Error: File not found")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line []string
	var days []Day
	cnt := 0

	for scanner.Scan() {
		if cnt != 0 {
			line = strings.Split(strings.ReplaceAll(scanner.Text(), "\"", ""), ",")
			days = append(days, Day{
				date:      line[0],
				fajr:      line[2],
				fajrSlutt: line[3],
				duhr:      line[4],
				asr:       line[5],
				maghrib:   line[8],
				isha:      line[9],
			})
		}
		cnt++
	}

	calendar := Calendar{days}
	fmt.Println(nextPrayer(calendar))
}

func nextPrayer(calendar Calendar) string {
	now := time.Now()

	for i := 0; i < 2; i++ {
		today, err := calendar.findDay(now.Format("02-01-2006"))
		if err != nil {
			return "Error: No data available"
		}

		prayerOrder := []struct {
			name string
			time string
		}{
			{"fajr", today.fajr},
			{"duhr", today.duhr},
			{"asr", today.asr},
			{"maghrib", today.maghrib},
			{"isha", today.isha},
		}

		for _, prayer := range prayerOrder {
			t, err := time.Parse("02-01-2006 15:04", fmt.Sprintf("%s %s", today.date, prayer.time))
			if err != nil {
				return fmt.Sprintf("Error parsing %s", prayer.name)
			}

			if now.Before(t) {
				if prayer.name == "fajr" {
					return fmt.Sprintf("%s: %s, soloppgang: %s", prayer.name, t.Format("15:04"), today.fajrSlutt)
				} else if prayer.name == "isha" {
					return fmt.Sprintf("%s: %s, midnatt: %s", prayer.name, t.Format("15:04"), Midnight(today.date, &calendar))
				} else {
					return fmt.Sprintf("%s: %s", prayer.name, t.Format("15:04"))
				}
			}
		}

		// Move to the next day if no prayer time remains today.
		now = now.Add(24 * time.Hour)
	}

	return "No upcoming prayer time"
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
