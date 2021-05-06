package main

import (
	"flag"
)

var (
	orderMessage string = "Order: Replace -> Retrieve -> Remove\n"
)

var (
	cmdlnBeatmap = flag.String("beatmap", "", `Specifies a certain beatmap.
	 If set to non "" - the program will search for given name and perform the magic
	 provided in settings if successful`)
	showOrder = flag.Bool("showOrder", false, "Prints an order in which functions are performed on a beatmap")
)
