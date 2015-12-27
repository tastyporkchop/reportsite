package main

import (
	"os"
	"tastyporkchop/reportsite/earthquake_poller/usgsatom"
	"testing"
)

func TestParseUSGSAtom(t *testing.T) {
	file, err := os.Open("test/all_hour.atom")
	if err != nil {
		t.Error(err)
	}
	defer file.Close()
	var data usgsatom.USGSFeed
	parseFeed(file, &data)
	if err != nil {
		t.Error(err)
	}

	point := data.Entry[0].Point
	if point != "38.7994995 -122.8085022" {
		t.Errorf("data.Entry[0].Point not set: expected:%s but found:%s", "38.7994995 -122.8085022", point)
	}
}
