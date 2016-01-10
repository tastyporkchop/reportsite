package quakeml

import (
	"encoding/xml"
	"github.com/davecgh/go-spew/spew"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	qmlFile, err := os.Open("../test/all_hour.quakeml")
	if err != nil {
		t.Error(err)
	}
	defer qmlFile.Close()

	var qml Q
	dec := xml.NewDecoder(qmlFile)
	err = dec.Decode(&qml)
	if err != nil {
		t.Error(err)
	}
	t.Log(spew.Sdump(qml))
}
