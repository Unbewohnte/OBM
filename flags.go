package main

import (
	"flag"
)

var (
	cmdlnBeatmap = flag.String("beatmap", "", `Specifies a certain beatmap.
	 If set to non "" - the program will search for given name and perform the magic
	 provided in settings if successful`)
)
