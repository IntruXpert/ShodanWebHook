package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/mattn/go-sqlite3"
)

var bot *tgbotapi.BotAPI
var channelID int64 = -1001831341573 // Replace with your own channel ID

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./shodan_data.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS shodan_data (
		trigger TEXT,
		ip TEXT,
		port INTEGER,
		isp TEXT,
		country_code TEXT,
		country_name TEXT,
		hostnames TEXT,
		asn TEXT,
		org TEXT,
		os TEXT,
		product TEXT,
		waf TEXT,
		title TEXT,
		status TEXT
	);`

	if _, err = db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}

	return db
}

func takeScreenshot(ip string, port int, org, product, title, hostnames, asn string) error {
	protocols := []string{"http", "https"}

	for _, protocol := range protocols {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.NoFirstRun,
			chromedp.NoDefaultBrowserCheck,
			chromedp.Headless,
		)

		actx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
		defer cancel()

		ctx, cancel := chromedp.NewContext(actx)
		defer cancel()

		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		var buf []byte
		if err := chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("%s://%s:%d", protocol, ip, port)),
			chromedp.CaptureScreenshot(&buf),
		); err != nil {
			log.Printf("Failed to take screenshot with %s: %v", protocol, err)
			continue
		}

		fileName := fmt.Sprintf("sc/%s_%s_%d.jpg", protocol, ip, port)
		if err := ioutil.WriteFile(fileName, buf, 0644); err != nil {
			log.Println("Failed to save screenshot:", err)
			return err
		}

		caption := fmt.Sprintf("%s://%s:%d \nOrg:[ %s ] \nProduct:[ %s ] \nTitle:[ %s ] \nHostnames:[ %s ] \nASN:[ %s ]", protocol, ip, port, org, product, title, hostnames, asn)
		photo := tgbotapi.NewPhotoUpload(channelID, fileName)
		photo.Caption = caption

		if _, err := bot.Send(photo); err != nil {
			log.Println("Failed to send photo:", err)
		}

		time.Sleep(3 * time.Second)
	}
	return nil
}

func handleShodanWebhook(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	defer r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Could not read incoming update:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		log.Println("Could not decode incoming update:", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	trigger, _ := data["trigger"].(string)
	ip, _ := data["ip_str"].(string)
	port, _ := data["port"].(float64)
	asn, _ := data["asn"].(string)
	isp, _ := data["isp"].(string)
	location, _ := data["location"].(map[string]interface{})
	countryCode, _ := location["country_code"].(string)
	countryName, _ := location["country_name"].(string)
	org, _ := data["org"].(string)
	os, _ := data["os"].(string)

	httpData, _ := data["http"].(map[string]interface{})
	title, _ := httpData["title"].(string)
	status, _ := httpData["status"].(string)
	waf, _ := httpData["waf"].(string)
	hostnames, _ := data["hostnames"].(string)
	product, _ := data["product"].(string)

	checkSQL := `SELECT ip, port FROM shodan_data WHERE ip = ? AND port = ?`
	var existingIP string
	var existingPort int
	err = db.QueryRow(checkSQL, ip, int(port)).Scan(&existingIP, &existingPort)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Failed to check existing data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if existingIP == ip && existingPort == int(port) {
		log.Printf("IP %s and port %d already exist. Skipping screenshot.\n", ip, port)
		return
	}

	insertSQL := `INSERT INTO shodan_data (trigger, ip, port, isp, country_code, country_name, hostnames, asn, org, os, product, waf, title, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err = db.Exec(insertSQL, trigger, ip, int(port), isp, countryCode, countryName, hostnames, asn, org, os, product, waf, title, status); err != nil {
		log.Println("Failed to insert data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := takeScreenshot(ip, int(port), org, product, title, hostnames, asn); err != nil {
		log.Println("Failed to take screenshot:", err)
	}
}

func main() {
	db := initDB()
	defer db.Close()

	var err error
	bot, err = tgbotapi.NewBotAPI("YOUR_TELEGRAM_BOT_API_TOKEN") // Replace with your bot token
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/updat1X73rj92", func(w http.ResponseWriter, r *http.Request) {
		handleShodanWebhook(w, r, db)
	})

	server := &http.Server{
		Addr: ":9080",
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	select {}
}
