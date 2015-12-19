package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	_ "github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	_ "golang.org/x/tools/blog/atom"
	"log"
	"net/http"
	"tastyporkchop/reportsite/earthquake_poller/usgsatom"
	"time"
)

var atomUrl string
var pollInterval time.Duration
var dbConStr string
var logFile string
var verbose bool

func main() {
	flag.StringVar(&atomUrl, "feedUrl", "http://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.atom", "feed url string")
	flag.DurationVar(&pollInterval, "pollInterval", time.Duration(5*time.Minute), "ms, s, or m")
	flag.StringVar(&dbConStr, "connStr", "user=angus database=reportsite", "")
	flag.StringVar(&logFile, "logFile", "", "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()

	// set up the log file

	// set up db
	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// poll
	poller := func(url string) (usgsatom.USGSFeed, error) {
		resp, err := http.Get(url)
		if err != nil {
			return usgsatom.USGSFeed{}, err
		}
		defer resp.Body.Close()
		decoder := xml.NewDecoder(resp.Body)
		var data usgsatom.USGSFeed
		err = decoder.Decode(&data)
		if err != nil {
			return usgsatom.USGSFeed{}, err
		}
		return data, err
	}
	feed, err := poller(atomUrl)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(spew.Sdump(feed))

	for i := range feed.Entry {
		err := processEntry(feed.Entry[i], db)
		if err != nil {
			log.Printf("Trouble processing entry: %s", err)
		}
	}
}

func processEntry(entry *usgsatom.USGSEntry, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO earthquake_event (event_id, title, updated, link, summary, summary_type, point, elevation) VALUES($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err := stmt.Exec(entry.ID, entry.Title, entry.Updated, entry.Link[0], entry.Summary.Body, entry.Summary.Type, entry.Point, entry.Elev)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
