# OBM (Osu!-background-manager)

## This utility will help you with replacement of Osu!`s beatmap backgrounds (and more in the future)

**Use at your own risk !**
There is no way to return removed original backgrounds unless you delete all beatmaps and reimport newly downloaded versions of them again.  

---

## Installation 

### From source
1. `git clone https://github.com/Unbewohnte/OBM.git` or download and unzip the archive
2. cd into the directory
3. `go build`

### From release
1. go to the [releases](https://github.com/Unbewohnte/OBM/releases) page
2. choose your OS and download the archive
3. cd to the location of the downloaded version
4. unzip (`7z x **archive_name**`) - for 7z archives 

---

## Usage

### First run 
1. The program will generate a settings.json file if it is not already in the directory when you run it
2. Paste your Osu! filepath in the "pathToOsu" field
3. Paste the filepath of the image in the "pathToimage" field. **ALL** beatmap`s backgrounds will be replaced with this image 
4. Additionally you can disable the "createBlackBackgoundImage" by replacing **true** with **false** or change the number of workers
5. Run the program once again

### After
1. Just run the utility again. If it found the settings file - it will perform the magic

---

## License
MIT License

---

If you have found this program useful, then consider to give this repository a ☆. It is not difficult for you, but means a lot for me 
