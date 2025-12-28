package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Playlist struct {
	items []Video
	xspf  string
}

const (
	XML_START = `<?xml version="1.0" encoding="UTF-8"?><playlist xmlns="http://xspf.org/ns/0/" xmlns:vlc="http://www.videolan.org/vlc/playlist/ns/0/" version="1"><title>Playlist</title>`
	XML_END   = `</playlist>`

	EXT_START = `<extension application="http://www.videolan.org/vlc/playlist/0">`
	EXT_END   = `</extension>`
)

func (pl *Playlist) Build() {

	if len(pl.items) == 0 {
		return
	}

	var s strings.Builder

	s.WriteString(XML_START)

	s.WriteString("<trackList>")

	for i, item := range pl.items {

		rawPath, _ := exec.Command("wslpath", "-w", item.path).Output()
		winPath := strings.TrimSpace(string(rawPath))

		uri := "file:///" + winPath

		s.WriteString("<track>")

		s.WriteString(fmt.Sprintf("<location>%s</location>", uri))

		s.WriteString(EXT_START)

		s.WriteString(fmt.Sprintf("<vlc:id>%d</vlc:id>", i))

		if item.start > 0 {
			s.WriteString(fmt.Sprintf("<vlc:option>start-time=%d</vlc:option>", item.start))
		}

		s.WriteString(EXT_END)

		s.WriteString("</track>")

	}

	s.WriteString("</trackList>")

	s.WriteString(EXT_START)

	for i := range pl.items {
		s.WriteString(fmt.Sprintf("<vlc:item tid=\"%d\"/>", i))
	}

	s.WriteString(EXT_END)

	s.WriteString(XML_END)

	pl.xspf = s.String()

}

func (pl Playlist) Save() error {

	return os.WriteFile(workingDir+"tis.xspf", []byte(pl.xspf), 0644)

}
