# OBM (Osu!-Background-Manager)

## This utility will help you with replacement, retrievement and removement of Osu!\`s beatmaps\` backgrounds

**Use at your own risk !**
There is no way to return removed original backgrounds unless you delete all beatmaps and reimport newly downloaded versions of them again.  

---

## Installation 

### From source (You`ll need [Go](https://golang.org/dl/) installed)
1. `git clone https://github.com/Unbewohnte/OBM.git` or download and unzip the archive
2. `cd` into the directory
3. `go build`

### From release
1. go to the [releases](https://github.com/Unbewohnte/OBM/releases) page
2. choose your OS and download the archive
3. `cd` to the location of the downloaded version
4. unzip (`7z x **archive_name**`) - for 7z archives

---

## Usage
To run - `./OBM` in terminal (on Unix) || `OBM` in command line (on Windows) (a simple double-click on exe will certainly work as well)

### First run 
1. The program will generate a settings.json file if it is not already in the directory when you run it
2. Paste your Osu! filepath in the "pathToOsu" field
3. Enable/Disable needed features, providing valid filepaths to them 
4. Additionally you can disable the "createBlackBackgoundImage" by replacing **true** with **false** or change the number of workers

### After
1. Start the utility again. If it has found the settings file - it will perform the magic according to provided rules


### Flags (starting from version 1.3.4)
Right now there is 2 arguments that you can specify before running the program - "beatmap" and "showOrder".
"-beatmap" flag takes a string; it will tell the program to do its work **ONLY** on beatmaps with specified name; others will be ignored.
The names of beatmaps in Osu! consist an id, artist and the name of the soundtrack, so you can
specify any name in the flag that will contain one of those parts.

"-showOrder" flag takes a boolean; if set to **true** - it will print an order in which the workers perform operations over each beatmap. Right now it`s just a helper flag.

#### Examples
1. `./OBM -beatmap=""` - the same as just `./OBM`. It will affect **all** of your beatmaps
2. `./OBM -beatmap="Demetori"` - this will search for beatmaps with names that contain "Demetori" and will work only with them
3. `./OBM -beatmap=Demetori` - the same as before, but without "" (look at 4 - 5 for nuances)
4. `./OBM -beatmap=raise my sword` - NOTE that this will make the program look only for word "raise", but not for the whole sequence
5. `./OBM -beatmap="raise my sword"` - this is the valid option for 4 (You need to use "" in case of a multi-word name)

The search is case-insensitive, so for example `./OBM -beatmap="Road of Resistance"` and `./OBM -beatmap="ROAD of rEsIsTaNcE"` will get you the same results

6. `./OBM -showOrder=true` - will print the order and exit
7. `./OBM -showOrder=true -beatmap="something here"` - will print the order and exit, just like in the previous one
---

## License
MIT License

---

If you have found this program useful, then consider to give this repository a ☆. It is not difficult for you, but means a lot for me 
