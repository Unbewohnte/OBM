package main

import (
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	// used as a flag if the program executed for "the first time"
	settingsFileExisted bool = false
)

const (
	settingsFilename string = "settings.json"
)

// struct for json settings` file contents
type Settings struct {
	OsuDir               string `json:"pathToOsu"`
	ReplacementImagePath string `json:"pathToimage"`
	CreateBlackBGImage   bool   `json:"createBlackBackgoundImage"`
}

// creates directory for logs and sets output to file
func setUpLogs() {
	logsDir := filepath.Join(".", "logs")
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(filepath.Join(logsDir, "logs.log"))
	log.SetOutput(file)
}

// creates "settings.json" and sets the flag
func createSettingsFile() {
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal("ERROR : Unable to read current directory")
	}
	for _, file := range files {
		if file.IsDir() == false {
			if file.Name() == settingsFilename {
				log.Println("Found settings file")
				settingsFileExisted = true
				return
			}
		}
	}
	file, err := os.Create("settings.json")
	if err != nil {
		log.Fatal("ERROR: Error creating settings file... : ", err)
	}
	settings := Settings{
		OsuDir:               "",
		ReplacementImagePath: "",
		CreateBlackBGImage:   true,
	}
	jsonEncodedSettings, err := json.MarshalIndent(settings, "", " ")
	if err != nil {
		log.Println("ERROR: Error creating settings file... : ", err)
	}
	file.Write(jsonEncodedSettings)

	file.Close()
	log.Println("Successfully created new settingsFile")
}

// filepath.Joins the main osu directory with its songs folder
func getSongsDir(osudir string) string {
	songsDir := filepath.Join(osudir, "Songs")

	stat, err := os.Stat(songsDir)
	if err != nil {
		log.Fatal("ERROR: Error reading path : ", err)
	}
	if !stat.IsDir() {
		log.Fatal("ERROR: Given osu! directory is not a directory")
	}

	return songsDir
}

// unmarshalls settings.json into struct
func getSettings() Settings {
	settingsFile, err := os.ReadFile(settingsFilename)
	if err != nil {
		log.Fatal("ERROR: Could not read settings file : ", err)
	}
	var settings Settings
	err = json.Unmarshal(settingsFile, &settings)
	if err != nil {
		log.Fatal("ERROR: Error unmarshalling json file : ", err)
	}
	return settings
}

// checks if given string contains ".osu"
func isBeatmap(filename string) bool {
	if len(filename) < 5 {
		return false
	}
	if filename[len(filename)-4:] == ".osu" {
		return true
	}
	return false
}

// parses .osu file and returns the filename of its background
func getBackgroundName(pathToOSUbeatmap string) string {
	beatmapBytes, err := os.ReadFile(pathToOSUbeatmap)
	if err != nil {
		log.Println("ERROR: Error reading beatmap file : ", err)
	}
	beatmapContents := string(beatmapBytes)

	eventsIndex := strings.Index(beatmapContents, "[Events]")
	if eventsIndex == -1 {
		return ""
	}
	breakPeriodsIndex := strings.Index(beatmapContents, "//Break Periods")
	if eventsIndex == -1 {
		return ""
	}
	beatmapBackground := strings.Split(beatmapContents[eventsIndex:breakPeriodsIndex], ",")[2]

	return beatmapBackground
}

// opens given files, copies one into another
func copyFile(src, dst string) {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Println("ERROR: Error copying files : ", err)
	}
}

// reads contents of given dir; searches for .osu files; parses them for background info;
// removes original background and replaces it with copied version of given image
func replaceBackgrounds(beatmapFolder, replacementPicPath string) {
	files, err := os.ReadDir(beatmapFolder)
	if err != nil {
		log.Fatal("ERROR: Wrong path : ", err)
	}
	for _, file := range files {
		filename := file.Name()

		if isBeatmap(filename) {

			beatmapBackgroundFilename := strings.Split(getBackgroundName(filepath.Join(beatmapFolder, filename)), "\"")[1]
			if beatmapBackgroundFilename == "" {
				continue
			}

			backgroundPath := filepath.Join(beatmapFolder, beatmapBackgroundFilename)
			log.Println(backgroundPath)

			// remove old background
			os.Remove(backgroundPath)

			// create new background file
			bgFile, err := os.Create(backgroundPath)
			if err != nil {
				log.Println("ERROR: Error creating new background file : ", err)
			}
			bgFile.Close()

			// copy the contents of a given image to the newly created bg file
			copyFile(replacementPicPath, backgroundPath)
		}

	}
}

// creates a complete black image file
func createBlackBG(width, height int) {
	bg, err := os.Create("blackBG.png")
	if err != nil {
		log.Println("ERROR: Error creating black background : ", err, "Continuing to run...")
	}
	image := image.NewRGBA(image.Rect(0, 0, width, height))
	bounds := image.Bounds()

	for y := 0; y < bounds.Max.Y; y++ {
		for x := 0; x < bounds.Max.X; x++ {
			image.Set(x, y, color.Black)
		}
	}
	png.Encode(bg, image)
	bg.Close()

	log.Println("Successfully created black background")
}

func init() {
	setUpLogs()
	createSettingsFile()
}

func main() {
	// settings file didn`t exist, created now
	if !settingsFileExisted {
		return
	}

	settings := getSettings()

	// process the given settings
	if settings.CreateBlackBGImage == true {
		createBlackBG(1920, 1080)
	}

	osuSongsDir := getSongsDir(settings.OsuDir)

	replacementImage := settings.ReplacementImagePath
	if replacementImage == "" || replacementImage == " " {
		log.Fatal("Image path not specified ! Specify `pathToimage` in settings file !")
	}

	// reading contents of `Songs` folder
	osuSongsDirContents, err := os.ReadDir(osuSongsDir)
	if err != nil {
		log.Fatal("ERROR: Error reading osu songs directory : ", err)
	}

	// storing all paths to each beatmap
	var songPaths []string
	for _, content := range osuSongsDirContents {
		if content.IsDir() {
			songPaths = append(songPaths, filepath.Join(osuSongsDir, content.Name()))
		}
	}
	log.Printf("Found %d song folders", len(songPaths))

	// replacing backgrounds for each beatmap
	for _, songPath := range songPaths {
		replaceBackgrounds(songPath, replacementImage)
	}

}
