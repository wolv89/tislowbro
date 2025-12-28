package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	slashDivider = " / "
	hashDivider  = " # "
)

func ParseFile(f *os.File) ([]Video, error) {

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	videos := make([]Video, 0)

	for scanner.Scan() {

		line := scanner.Text()

		if len(line) == 0 {
			continue
		}

		video, err := ParseVideo(line)
		if err != nil {
			fmt.Println("parse error:", err.Error())
			continue
		}

		videos = append(videos, video)

	}

	if len(videos) == 0 {
		return nil, fmt.Errorf("unable to parse any videos from given file")
	}

	return videos, nil

}

/*
 * Read video from line like:
 * dir / video-name
 * [or]
 * dir:sub / video-name # optional extras
 * The spacing is important, we check for the slash/hash with spacing
 */
func ParseVideo(line string) (Video, error) {

	v := Video{}

	if line == "" {
		return v, errors.New("unable to parse empty string")
	}

	slash := strings.Index(line, slashDivider)
	if slash < 1 {
		return v, errors.New("no directory given")
	}

	dir := strings.TrimSpace(line[:slash])
	if len(dir) == 0 {
		return v, errors.New("no directory given")
	}

	end := len(line)

	hash := strings.LastIndex(line, hashDivider)
	if hash != -1 {
		v.start = ParseParams(line, hash)
		end = hash
	}

	v.name = line[slash+len(slashDivider) : end]
	v.dir = dir

	err := v.Find()
	if err != nil {
		return v, err
	}

	return v, nil

}

// Currently only looking for one thing: start time
// Longer term this might need to return a struct of options
func ParseParams(line string, h int) uint {

	params := strings.TrimSpace(line[h+len(hashDivider):])
	if len(params) == 0 {
		return 0
	}

	last := strings.LastIndexByte(params, ' ')
	stamp := params[last+1:]

	if len(stamp) == 0 {
		return 0
	}

	var mins, secs uint
	fmt.Sscanf(stamp, "%d:%d", &mins, &secs)

	return (mins * 60) + secs

}
