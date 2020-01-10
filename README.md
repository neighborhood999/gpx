# gpx

A simple gpx parser (support Strava currently) for running and written in Go.

## Installation

```sh
$ make dep
```

## Tests

```sh
$ make test
```

## Usage

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/neighborhood999/gpx"
)

func main() {
	f, err := os.Open("running.gpx")

	if err != nil {
		fmt.Println(err)
	}

	defer f.Close()

	b, _ := ioutil.ReadAll(f)
	g, _ := gpx.ReadGPX(bytes.NewReader(b))

	fmt.Println(g.Distance()) // Get the total running distance
	fmt.Println(g.Duration()) // Get the total running duration
	fmt.Println(g.PaceInKM()) // Get the running pace(km/min)
}
```

## LICENSE

MIT © [Peng Jie](https://github.com/neighborhood999/)
