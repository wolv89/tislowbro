package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	flagFile string
	flagMode int

	workingDir string

	dirCache map[string][]os.DirEntry
)

const (
	MODE_STANDARD int = iota
	MODE_DEVELOPMENT
)

func init() {

	flag.StringVar(&flagFile, "file", "", "Path to file to use, relative or absolute")
	flag.StringVar(&flagFile, "f", "", "Path to file to use, relative or absolute")

	flag.IntVar(&flagMode, "mode", 0, "Running mode - 0 is standard, 1 is development")
	flag.IntVar(&flagMode, "m", 0, "Running mode - 0 is standard, 1 is development")

}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	vlc_addr := os.Getenv("VLC_ADDRESS")
	vlc_pass := os.Getenv("VLC_PASSWORD")

	if len(vlc_addr) == 0 {
		log.Fatal("missing VLC address in env")
	}
	if len(vlc_pass) == 0 {
		log.Fatal("missing VLC password in env")
	}

	workingDir = os.Getenv("WORK_DIR")

	if len(workingDir) == 0 {
		workingDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}

	player := NewPlayer(vlc_addr, vlc_pass)
	if err := player.IsRunning(); err != nil {
		fmt.Println("aborted launch, VLC player does not appear to be running, or is not configured for HTTP connections")
		log.Fatal(err.Error())
	}

	if err := player.Reset(); err != nil {
		fmt.Println("aborted launch, unable to reset VLC play list")
		log.Fatal(err.Error())
	}

	flag.Parse()

	if len(flagFile) == 0 {
		log.Fatal("file name/path required")
	}

	f, err := os.Open(flagFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	workingDir += "/"
	dirCache = make(map[string][]os.DirEntry)

	videos, err := ParseFile(f)
	if err != nil {
		log.Fatal(err)
	}

	playlist := Playlist{
		items: videos,
	}

	playlist.Build()
	err = playlist.Save()
	if err != nil {
		log.Fatal(err)
	}

	player.Queue(playlist)

}
