package main

import (
	"adult-short-videos/internal/config"
	"adult-short-videos/internal/service/video/model"
	"encoding/json"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const maxPages = 1

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64; rv:125.0) Gecko/20100101 Firefox/125.0",
}

func main() {
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	db, err := gorm.Open(postgres.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	client := &http.Client{Timeout: 30 * time.Second}

	totalInserted := 0
	emptyPages := 0

	for page := 1; page <= maxPages; page++ {
		htmlStr, err := fetchPage(client, page, rng)
		if err != nil {
			log.Printf("第 %d 页请求失败: %v", page, err)
			emptyPages++
			if emptyPages >= 3 {
				break
			}
			continue
		}

		inserted := parseAndUpsert(db, client, htmlStr)
		log.Printf("第 %d 页: 写入 %d 条", page, inserted)
		totalInserted += inserted

		if inserted == 0 {
			emptyPages++
			if emptyPages >= 3 {
				log.Println("连续 3 页无数据，停止采集")
				break
			}
		} else {
			emptyPages = 0
		}

		// 随机延迟 1.5 ~ 3.5 秒
		time.Sleep(time.Duration(1500+rng.Intn(2000)) * time.Millisecond)
	}

	log.Printf("采集完成，共写入 %d 条视频", totalInserted)
}

// fetchPage 返回页面的 HTML 字符串（page1 是完整 HTML，page2+ 是 JSON.data.html）
func fetchPage(client *http.Client, page int, rng *rand.Rand) (string, error) {
	url := "https://123av.fun/"
	if page > 1 {
		url = fmt.Sprintf("https://123av.fun/api/page-%d", page)
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgents[rng.Intn(len(userAgents))])
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://123av.fun/")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if page == 1 {
		return string(body), nil
	}

	// page 2+: {"message":"...","data":{"html":"..."}}
	var result struct {
		Data struct {
			HTML string `json:"html"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("json parse: %w", err)
	}
	return result.Data.HTML, nil
}

// detectPortrait 通过下载封面图头部字节判断是否竖屏（height > width）
func detectPortrait(client *http.Client, coverURL string) bool {
	if coverURL == "" {
		return false
	}
	req, err := http.NewRequest("GET", coverURL, nil)
	if err != nil {
		return false
	}
	req.Header.Set("User-Agent", userAgents[0])
	req.Header.Set("Referer", "https://123av.fun/")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	cfg, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return false
	}
	return cfg.Height > cfg.Width
}

// parseAndUpsert 解析 HTML 片段，upsert 到 DB，返回实际写入数量
func parseAndUpsert(db *gorm.DB, client *http.Client, htmlStr string) int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		log.Printf("goquery parse error: %v", err)
		return 0
	}

	inserted := 0
	doc.Find("a.video-card, a[data-src]").Each(func(_ int, s *goquery.Selection) {
		remoteID := strings.TrimSpace(s.AttrOr("data-id", ""))
		if remoteID == "" {
			return
		}

		title := strings.TrimSpace(s.AttrOr("data-videotitle", ""))
		if title == "" {
			return
		}
		title = strings.ReplaceAll(title, "123av.fun", "蜜臀69")
		coverURL := strings.TrimSpace(s.AttrOr("data-poster", ""))
		sourceURL := strings.TrimSpace(s.AttrOr("data-src", ""))
		if sourceURL == "" {
			return
		}
		durationSec := parseInt32(s.AttrOr("data-duration", ""))
		playCount := parseInt64(s.AttrOr("data-playcount", ""))

		isPortrait := detectPortrait(client, coverURL)

		// 尝试多个可能的作者/频道字段
		author := strings.TrimSpace(s.AttrOr("data-author",
			s.AttrOr("data-channel",
				s.AttrOr("data-publisher",
					s.AttrOr("data-uploader", "")))))
		// 也尝试子元素文本
		if author == "" {
			author = strings.TrimSpace(s.Find(".author, .channel, .publisher, .uploader").First().Text())
		}

		v := model.Video{
			Title:       title,
			CoverURL:    coverURL,
			Duration:    durationSec,
			StorageType: "cold",
			SourceURL:   sourceURL,
			PlayCount:   playCount,
			IsPortrait:  isPortrait,
			Author:      author,
			Status:      1,
			RemoteId:    remoteID,
			PublishedAt: time.Now(),
		}

		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "remote_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"title", "cover_url", "source_url", "play_count", "duration", "is_portrait", "author"}),
		}).Create(&v)

		if result.Error != nil {
			log.Printf("upsert error (remote_id=%s): %v", remoteID, result.Error)
			return
		}
		inserted++
	})
	return inserted
}

func parseInt32(s string) int32 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 32)
	return int32(v)
}

func parseInt64(s string) int64 {
	v, _ := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	return v
}
