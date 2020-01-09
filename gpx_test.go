package gpx

import (
	"bytes"
	"io/ioutil"
	"log"
	"math"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testGPX = "_data/strava-running-sample.gpx"
)

func openGPX(path string) *bytes.Reader {
	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)

	if err != nil {
		log.Fatal(err)
	}

	return bytes.NewReader(b)
}

func TestRead(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	assert.Equal(t, "StravaGPX", gpx.Creator)
	assert.Equal(t, "1.1", gpx.Version)
	assert.IsType(t, &MetaData{}, gpx.Metadata)
	assert.Len(t, gpx.Tracks, 1)
}

func TestWayPointTime(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	firstTrakPointTime := gpx.Tracks[0].TrackSegments[0].TrackPoint[0].Time()
	expectedTime := time.Date(2019, 10, 26, 21, 21, 11, 0, time.UTC)

	assert.Equal(t, expectedTime, firstTrakPointTime)
}

func TestDuration(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	assert.Equal(t, 34.0, gpx.Duration())
}

func TestZeroDuration(t *testing.T) {
	b := openGPX("_data/zero-duration.gpx")
	gpx, _ := ReadGPX(b)

	assert.Equal(t, 0.0, gpx.Duration())
}

func TestTwoPointDistance(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	start := gpx.Tracks[0].TrackSegments[0].TrackPoint[0]
	end := gpx.Tracks[0].TrackSegments[0].TrackPoint[5]

	assert.Less(t, float64(0.02), start.Distance(&end))
}

func TestGPXDistance(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	assert.Less(t, float64(0.1), gpx.Distance())
}

func TestPaceInKM(t *testing.T) {
	b := openGPX(testGPX)
	gpx, _ := ReadGPX(b)

	p := gpx.PaceInKM()

	assert.Equal(t, &Pace{4, 49}, p)
}

func TestToRadians(t *testing.T) {
	assert.Equal(t, math.Pi, toRadians(180))
}
