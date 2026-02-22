package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

const (
	baseURL       = "https://mypage.neis-gym.com/reserve/schedule/68/201/"
	schoolClassID = "108"
	weeksToCheck  = 4
)

// AvailableSlot represents an available reservation slot.
type AvailableSlot struct {
	Date      string // e.g., "3月22日(日)"
	Time      string // e.g., "15:00～16:00"
	ClassName string // e.g., "幼児クラス"
}

func (s AvailableSlot) String() string {
	return fmt.Sprintf("【空きあり】\n日付: %s\n時間: %s\nクラス: %s\n", s.Date, s.Time, s.ClassName)
}

// ScrapeAvailableSlots scrapes the website and returns a list of available slots.
func ScrapeAvailableSlots() ([]AvailableSlot, error) {
	var allAvailableSlots []AvailableSlot

	// 実行日から4週間後まで、7日おきに計4回のクエリを発行
	for i := 0; i < weeksToCheck; i++ {
		dateFrom := time.Now().AddDate(0, 0, 7*i).Format("2006-01-02")
		targetURL := fmt.Sprintf("%s?school_class_id=%s&date_from=%s", baseURL, schoolClassID, dateFrom)

		log.Printf("Scraping URL: %s", targetURL)

		slots, err := parseSchedulePage(targetURL)
		if err != nil {
			// 1つのページの取得に失敗しても処理を続行するが、エラーはログに出力
			log.Printf("WARN: Failed to parse page for date %s: %v", dateFrom, err)
			continue
		}
		allAvailableSlots = append(allAvailableSlots, slots...)
	}

	return allAvailableSlots, nil
}

// parseSchedulePage fetches and parses a single schedule page.
func parseSchedulePage(url string) ([]AvailableSlot, error) {
	// create context with suppressed logging
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.NoSandbox,
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(func(string, ...interface{}) {}))
	defer cancel()

	// navigate and get DOM content
	var html string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// wait for the page to load
		chromedp.Sleep(5*time.Second),
		chromedp.OuterHTML("html", &html),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve page content: %w", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var availableSlots []AvailableSlot

	// 日付ラベルを取得
	var dayLabels []string
	doc.Find(".week_label.days.fs_3 > div").Each(func(i int, s *goquery.Selection) {
		label := strings.TrimSpace(s.Text())
		dayLabels = append(dayLabels, label)
	})

	doc.Find(".lesson_container .days .day").Each(func(i int, daySel *goquery.Selection) {
		dateLabel := ""
		if i < len(dayLabels) {
			dateLabel = dayLabels[i]
		}
		// 曜日が月・土・日のみ対象
		if !(strings.Contains(dateLabel, "(土)") || strings.Contains(dateLabel, "(日)")) {
			return
		}
		// 各部屋の各レッスンを取得
		daySel.Find(".rooms .room .lessons .d_lesson").Each(func(j int, lessonSel *goquery.Selection) {
			if lessonSel.HasClass("full") {
				return
			}
			timeRange := strings.TrimSpace(lessonSel.Find(".contents .fs_2.mb_text").Text())
			className := strings.TrimSpace(lessonSel.Find(".school-class-lesson-schedule-box .schedule-label span").Text())
			if className != "" && timeRange != "" && dateLabel != "" {
				slot := AvailableSlot{
					Date:      dateLabel,
					Time:      timeRange,
					ClassName: className,
				}
				availableSlots = append(availableSlots, slot)
			}
		})
	})

	return availableSlots, nil
}
