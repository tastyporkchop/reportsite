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
	"os"
	"os/signal"
	"syscall"
	"tastyporkchop/reportsite/earthquake_poller/usgsatom"
	"time"
)

var atomUrl string
var pollInterval time.Duration
var dbConStr string
var logFile string
var verbose bool

func main() {
	log.Print("startup")
	log.Print("parsing args")
	flag.StringVar(&atomUrl, "feedUrl", "http://earthquake.usgs.gov/earthquakes/feed/v1.0/summary/all_hour.atom", "feed url string")
	flag.DurationVar(&pollInterval, "pollInterval", time.Duration(5*time.Minute), "ms, s, or m")
	flag.StringVar(&dbConStr, "connStr", "user=angus database=reportsite", "")
	flag.StringVar(&logFile, "logFile", "", "")
	flag.BoolVar(&verbose, "v", false, "")
	flag.Parse()

	// set up the log file

	// set up db
	log.Print("setting up database connection")
	db, err := sql.Open("postgres", dbConStr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { db.Close(); log.Print("database connections closed") }()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// set up the poller
	log.Print("setting up poller")
	pollTicker := time.NewTicker(pollInterval)
	defer func() {
		log.Print("stopping poller")
		pollTicker.Stop()
		log.Print("poller stopped")
	}()

	//
	go func() {
		for range pollTicker.C {
			log.Print("polling...")
			feed, err := poll(atomUrl)
			if err != nil {
				log.Printf("Trouble polling url:%s message:%s", atomUrl, err)
				continue
			}

			tx, err := db.Begin()
			if err != nil {
				log.Printf("Trouble beginning transaction:%s", err)
				continue
			}
			defer tx.Rollback()

			log.Printf("processing %d entries", len(feed.Entry))
			for i := range feed.Entry {
				err := processEntry(feed.Entry[i], tx)
				if err != nil {
					log.Printf("Trouble processing entry: %s", err)
				}
			}

			log.Print("commiting entries")
			err = tx.Commit()
			if err != nil {
				log.Printf("Trouble committing transaction: %s", err)
			}
		}
	}()

	// wait for interrupt
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		switch sig {
		case syscall.SIGTERM:
			log.Printf("Received SIGTERM:%s", sig)
			done <- true
		default:
			log.Printf("We don't do anything with signal: %s", sig)
		}
	}()

	// block until done
	<-done

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
	log.Print("insterting entry")
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
