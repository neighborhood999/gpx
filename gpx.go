// GPX 1.1 Schema Documentation: https://www.topografix.com/GPX/1/1/

package gpx

import (
	"encoding/xml"
	"io"
	"math"
	"time"

	"golang.org/x/net/html/charset"
)

// EARTHRADIUS is the Earth's radius 6,371km.
// ref: https://en.wikipedia.org/wiki/Earth_radius
const EARTHRADIUS = 6371

// GPX is the representation gpxType.
type GPX struct {
	XMLName  string    `xml:"gpx"`
	Creator  string    `xml:"creator,attr,omitempty"`
	Version  string    `xml:"version,attr,omitempty"`
	Metadata *MetaData `xml:"metadata,omitempty"`
	Tracks   []Track   `xml:"trk,omitempty"`
}

// MetaData is the information about the GPX file, author,
// and copyright restrictions goes in the metadata section.
type MetaData struct {
	XMLName   xml.Name `xml:"metadata"`
	Timestamp string   `xml:"time,omitempty"`
}

// Link is an external resource (Web page, digital photo, video clip, etc)
// with additional information.
type Link struct {
	XMLName xml.Name `xml:"link"`
	URL     string   `xml:"href,attr,omitempty"`
	Text    string   `xml:"text,omitempty"`
	Type    string   `xml:"type,omitempty"`
}

// Track is the representation trk - an ordered list of points describing a path.
type Track struct {
	XMLName       xml.Name       `xml:"trk"`
	Name          string         `xml:"name,omitempty"`
	Comment       string         `xml:"cmt,omitempty"`
	Description   string         `xml:"desc,omitempty"`
	Source        string         `xml:"src,omitempty"`
	Links         []Link         `xml:"link,omitempty"`
	Number        int            `xml:"number,omitempty"`
	Type          string         `xml:"type,omitempty"`
	Extensions    *Extensions    `xml:"extensions,omitempty"`
	TrackSegments []TrackSegment `xml:"trkseg,omitempty"`
}

// Extensions is the representation extension.
type Extensions struct {
	XML []byte `xml:",innerxml"`
}

// TrackSegment holds a list of TrackPoint which are logically connected in order.
type TrackSegment struct {
	XMLName    xml.Name    `xml:"trkseg"`
	TrackPoint []WayPoint  `xml:"trkpt"`
	Extensions *Extensions `xml:"extensions,omitempty"`
}

// WayPoint is a point of interest, or named feature on a map.
type WayPoint struct {
	XMLName                       xml.Name              `xml:"trkpt"`
	Latitude                      float64               `xml:"lat,attr"`
	Longitude                     float64               `xml:"lon,attr"`
	Elevation                     float64               `xml:"ele,omitempty"`
	Timestamp                     string                `xml:"time,omitempty"`
	MagneticVariation             Degrees               `xml:"magvar,omitempty"`
	GeoIDHeight                   float64               `xml:"geoidheight,omitempty"`
	Name                          string                `xml:"name,omitempty"`
	Comment                       string                `xml:"cmt,omitempty"`
	Description                   string                `xml:"desc,omitempty"`
	Source                        string                `xml:"src,omitempty"`
	Links                         []Link                `xml:"link,omitempty"`
	Symbol                        string                `xml:"sym,omitempty"`
	Type                          string                `xml:"type,omitempty"`
	Fix                           Fix                   `xml:"fix,omitempty"`
	Sat                           int                   `xml:"sat,omitempty"`
	HorizontalDilutionOfPrecision float64               `xml:"hdop,omitempty"`
	VerticalDilutionOfPrecision   float64               `xml:"vdop,omitempty"`
	PositionDilutionOfPrecision   float64               `xml:"pdop,omitempty"`
	AgeOfGpsData                  float64               `xml:"ageofgpsdata,omitempty"`
	DifferentialGPSID             DGPSStation           `xml:"dgpsid,omitempty"`
	Extensions                    *TrackPointExtensions `xml:"extensions,omitempty"`
}

// TrackPointExtensions extend GPX by adding your own elements from another schema
type TrackPointExtensions struct {
	XMLName              xml.Name             `xml:"extensions"`
	TrackPointExtensions *TrackPointExtension `xml:"TrackPointExtension,omitempty"`
}

// TrackPointExtension tracks temperature, heart rate and cadence specific to devices
type TrackPointExtension struct {
	XMLName      xml.Name `xml:"TrackPointExtension"`
	Temperature  float64  `xml:"atemp,omitempty"`
	WTemperature float64  `xml:"wtemp,omitempty"`
	Depth        float64  `xml:"depth,omitempty"`
	HeartRate    int      `xml:"hr,omitempty"`
	Cadence      int      `xml:"cad,omitempty"`
}

// Degrees is used for bearing, heading, course. Units are decimal degrees, true (not magnetic). (0.0 <= value < 360.0)
type Degrees float64

// Fix is the representation type of GPS fix. none means GPS had no fix.
// To signify "the fix info is unknown, leave out fixType entirely.
// pps = military signal used (value comes from list: {'none'|'2d'|'3d'|'dgps'|'pps'})
type Fix string

// DGPSStation is the representation a differential GPS station. (0 <= value <= 1023)
type DGPSStation int

// Pace is the representation a running pace.
type Pace struct {
	Minutes int
	Seconds int
}

// Point is the representation a point of latitude and longitude
type Point struct {
	Latitude  float64
	Longitude float64
}

// ReadGPX is a GPX reader and return a GPX object and error.
func ReadGPX(r io.Reader) (*GPX, error) {
	gpx := &GPX{}

	// ref: https://stackoverflow.com/questions/6002619/unmarshal-an-iso-8859-1-xml-input-in-go
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReaderLabel
	err := d.Decode(gpx)

	return gpx, err
}

// Time returns TrackPoint timestamp as Time
func (w *WayPoint) Time() time.Time {
	t, err := time.Parse(time.RFC3339, w.Timestamp)

	if err != nil {
		return time.Time{}
	}

	return t
}

// Distance returns two point distance.
// ref: https://www.movable-type.co.uk/scripts/latlong.html
func (w *WayPoint) Distance(w2 *WayPoint) float64 {
	lat1 := toRadians(w.Latitude)
	lat2 := toRadians(w2.Latitude)
	distanceLat := toRadians(w2.Latitude - w.Latitude)
	distanceLon := toRadians(w2.Longitude - w.Longitude)

	a := math.Sin(distanceLat/2)*math.Sin(distanceLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(distanceLon/2)*math.Sin(distanceLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EARTHRADIUS * c
}

// Duration returns the duration of all tracks in a GPX in seconds.
func (g *GPX) Duration() float64 {
	trackPoints := g.Tracks[0].TrackSegments[0].TrackPoint

	start := trackPoints[0].Time()
	end := trackPoints[len(trackPoints)-1].Time()

	if end.Equal(start) || end.Before(start) {
		return 0.0
	}

	duration := end.Sub(start)

	return duration.Seconds()
}

// Distance returns total distance
func (g *GPX) Distance() float64 {
	var totalDistance float64
	trackPoints := g.Tracks[0].TrackSegments[0].TrackPoint

	for i := 1; i < len(trackPoints); i++ {
		totalDistance += trackPoints[i-1].Distance(&trackPoints[i])
	}

	return totalDistance
}

// PaceInKM returns running pace in kilometers.
func (g *GPX) PaceInKM() *Pace {
	paceInKM := int(g.Duration() / g.Distance())
	minutesPaceInKm := int(paceInKM / 60)
	secondsPaceInKm := paceInKM % 60

	return &Pace{minutesPaceInKm, secondsPaceInKm}
}

// PaceInMile returns running pace in miles.
func (g *GPX) PaceInMile() *Pace {
	paceInKM := int(g.Duration() / g.Distance() / 1.609344)
	minutesPaceInKm := int(paceInKM / 60)
	secondsPaceInKm := paceInKM % 60

	return &Pace{minutesPaceInKm, secondsPaceInKm}
}

// Elevations returns all the track point elevation.
func (g *GPX) Elevations() []float64 {
	trackPoints := g.Tracks[0].TrackSegments[0].TrackPoint
	elevations := make([]float64, len(trackPoints))

	for i := range trackPoints {
		elevations[i] = trackPoints[i].Elevation
	}

	return elevations
}

// MinAndMixElevation returns min and mix elevation.
func (g *GPX) MinAndMixElevation() (float64, float64) {
	e := g.Elevations()
	minElevation := e[0]
	maxElevation := e[0]

	for _, value := range e {
		if value < minElevation {
			minElevation = value
		}

		if value > maxElevation {
			maxElevation = value
		}
	}

	return minElevation, maxElevation
}

// GetCoordinates return all track points latitude and longitude.
func (g *GPX) GetCoordinates() []Point {
	trackPoints := g.Tracks[0].TrackSegments[0].TrackPoint
	coordinates := make([]Point, len(trackPoints))

	for i, track := range trackPoints {
		coordinates[i] = Point{Longitude: track.Longitude, Latitude: track.Latitude}
	}

	return coordinates
}

// toRadians converts degrees to radians.
func toRadians(degree float64) float64 {
	return degree * math.Pi / 180.0
}
