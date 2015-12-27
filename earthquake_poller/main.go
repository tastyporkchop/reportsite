package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	_ "github.com/davecgh/go-spew/spew"
	_ "github.com/lib/pq"
	_ "golang.org/x/tools/blog/atom"
	"io"
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
	feed, err := poll(atomUrl)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(spew.Sdump(feed))

	//
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()

	for i := range feed.Entry {
		err := processEntry(feed.Entry[i], tx)
		if err != nil {
			log.Printf("Trouble processing entry: %s", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// poll
func poll(url string) (usgsatom.USGSFeed, error) {
	resp, err := http.Get(url)
	if err != nil {
		return usgsatom.USGSFeed{}, err
	}
	defer resp.Body.Close()
	var data usgsatom.USGSFeed
	err = parseFeed(resp.Body, &data)
	if err != nil {
		return usgsatom.USGSFeed{}, err
	}
	return data, nil
}

//
func parseFeed(r io.Reader, msg *usgsatom.USGSFeed) error {
	decoder := xml.NewDecoder(r)
	err := decoder.Decode(msg)
	if err != nil {
		return err
	}
	return nil
}

//
func entryExists(entry *usgsatom.USGSEntry, tx *sql.Tx) bool {
	// check if event exists
	stmt, err := tx.Prepare("SELECT COUNT(*) from earthquake_event where event_id=$1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	var idcount int
	err = stmt.QueryRow(entry.ID).Scan(&idcount)
	if err != nil {
		log.Fatal(err)
	}
	return idcount > 0
}

//
func insertEntry(entry *usgsatom.USGSEntry, tx *sql.Tx) error {
	// insert the record
	stmt, err := tx.Prepare("INSERT INTO earthquake_event (event_id, title, updated, link, summary, summary_type, point, elevation) VALUES($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(entry.ID, entry.Title, string(entry.Updated), entry.Link[0].Href, entry.Summary.Body, entry.Summary.Type, entry.Point, entry.Elev)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

//
func processEntry(entry *usgsatom.USGSEntry, tx *sql.Tx) error {

	if entryExists(entry, tx) {
		return nil
	}
	return insertEntry(entry, tx)
}
