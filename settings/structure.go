package settings

// the idea behind `Enabled` field is that if it`s not true - then
// we treat the path below as "" (blank string) , which workers will just ignore
// (therefore will not perform the replacement or retrievement)

type backgroundReplacement struct {
	Enabled              bool   `json:"enabled"`
	ReplacementImagePath string `json:"pathToimage"`
}

type backgroundRetrievement struct {
	Enabled          bool   `json:"enabled"`
	RetrievementPath string `json:"retrievementPath"`
}

type creatingBG struct {
	Enabled bool `json:"enabled"`
	Width   uint `json:"width"`
	Height  uint `json:"height"`
}

// struct for json settings` file contents
type Settings struct {
	OsuDir                 string                 `json:"pathToOsu"`
	BackgroundReplacement  backgroundReplacement  `json:"backgroundReplacement"`
	BackgroundRetrievement backgroundRetrievement `json:"backgroundRetrievement"`
	CreateBlackBGImage     creatingBG             `json:"createBlackBackgoundImage"`
	Workers                int                    `json:"workers"`
}
