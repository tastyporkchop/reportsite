package quakeml

import (
	"time"
)

type Q struct {
	Events []Event `xml:"eventParameters>event"`
}

type Event struct {
	Datasource           string       `xml:"datasource,attr"`
	Eventsource          string       `xml:"eventsource,attr"`
	Eventid              string       `xml:"eventid,attr"`
	PublicID             string       `xml:"publicID,attr"`
	Description          Description  `xml:"description"`
	Origin               Origin       `xml:"origin"`
	Magnitude            Magnitude    `xml:"magnitude"`
	PreferredOriginID    string       `xml:"preferredOriginID"`
	PreferredMagnitudeID string       `xml:"preferredMagnitudeID"`
	Type                 string       `xml:"type"`
	CreationInfo         CreationInfo `xml:"creationInfo"`
}

type Description struct {
	Type string `xml:"type"`
	Text string `xml:"text"`
}

type Origin struct {
	Datasource        string            `xml:"datasource,attr"`
	Dataid            string            `xml:"dataid,attr"`
	Eventsource       string            `xml:"eventsource,attr"`
	Eventid           string            `xml:"eventid,attr"`
	PublicID          string            `xml:"publicID,attr"`
	Time              TimeStr           `xml:"time>value"`
	Longitude         float64           `xml:"longitude>value"`
	Latitude          float64           `xml:"latitude>value"`
	Depth             Depth             `xml:"depth"`
	OriginUncertainty OriginUncertainty `xml:"originUncertainty"`
	Quality           Quality           `xml:"quality"`
	EvaluationMode    string            `xml:"evaluationMode"`
	CreationInfo      CreationInfo      `xml:"creationInfo"`
}

type Depth struct {
	Value       int `xml:"value"`
	Uncertainty int `xml:"uncertainty"`
}

type OriginUncertainty struct {
	HorizontalUncertainty int    `xml:"horizontalUncertainty"`
	PreferredDescription  string `xml:"preferredDescription"`
}

type Quality struct {
	UsedPhaseCount   int     `xml:"usedPhaseCount"`
	UsedStationCount int     `xml:"usedStationCount"`
	StandardError    float64 `xml:"standardError"`
	AzimuthalGap     float64 `xml:"azimuthalGap"`
	MinimumDistance  float64 `xml:"minimumDistance"`
}

type CreationInfo struct {
	AgencyID     string  `xml:"agencyID"`
	CreationTime TimeStr `xml:"creationTime"`
	Version      string  `xml:"version"`
}

type Magnitude struct {
	Mag            float64      `xml:"mag>value"`
	Uncertainty    float64      `xml:"mag>uncertainty"`
	Type           string       `xml:"type"`
	StationCount   int          `xml:"stationCount"`
	OriginID       string       `xml:"originID"`
	EvaluationMode string       `xml:"evaluationMode"`
	CreationInfo   CreationInfo `xml:"creationInfo"`
}

type TimeStr string

func Time(t time.Time) TimeStr {
	//return TimeStr(t.Format("2006-01-02T15:04:05-07:00"))
	//2006-01-02T15:04:05.999999999Z07:00
	return TimeStr(t.Format(time.RFC3339Nano))
}
