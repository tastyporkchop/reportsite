package geojson

type FeatureCollection struct {
	Type     string
	BBox     BBox
	Features []Feature
}

type BBox struct {
	Minimum_longitude float64
	Minimum_latitude  float64
	Minimum_depth     float64
	Maximum_longitude float64
	Maximum_latitude  float64
	Maximum_depth     float64
}

type Feature struct {
	Type       string
	Geometry   Geometry
	Id         string
	Properties map[string]string
}

type Properties struct {
	Type    string
	Mag     float64
	Place   string
	Time    int
	Updated int
	Tz      int
	Url     string
	Detail  string
	Felt    int
	Cdi     float64
	Mmi     float64
	Alert   string
	Status  string
	Tsunami int
	Sig     int
	Net     string
	Code    string
	Ids     string
	Sources string
	Types   string
	Nst     int
	Dmin    float64
	Rms     float64
	Gap     float64
	MagType string
}
