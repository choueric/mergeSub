package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputSubsStr  string
	timeOffsetStr string
	formatStr     string

	inputSubs  []string
	outputSub  string
	timeOffset []string
	lineEnd    string

	debug bool
)

func parseOptions() error {
	if inputSubsStr == "" || outputSub == "" || timeOffsetStr == "" {
		return errors.New("invalid parameters")
	}
	inputSubs = strings.Split(inputSubsStr, ";")
	timeOffset = strings.Split(timeOffsetStr, ";")

	switch formatStr {
	case "unix":
		lineEnd = "\n"
	case "dos":
		lineEnd = "\r\n"
	default:
		return errors.New("invalid format string.")
	}

	fmt.Println("- input sub files:")
	for i, v := range inputSubs {
		fmt.Println(" ", i, ":", v)
	}
	fmt.Println("- output sub file:\n ", outputSub)
	fmt.Println("- time offset:")
	for i, v := range timeOffset {
		fmt.Println(" ", i, ":", v)
	}
	fmt.Println("- output format:", formatStr)

	return nil
}

func main() {
	flag.StringVar(&inputSubsStr, "i", "", "input subtitle files, separated by ';'.")
	flag.StringVar(&outputSub, "o", "", "ouput subtitle file.")
	flag.StringVar(&timeOffsetStr, "t", "", "time offset, separated by ';'.")
	flag.StringVar(&formatStr, "f", "unix", "output file format, unix or dos.")
	flag.BoolVar(&debug, "d", false, "enable debug message.")
	flag.Parse()

	err := parseOptions()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	wait := make(chan int, len(inputSubs))

	srtFiles := []*SrtFile{}
	for _, v := range inputSubs {
		srt := &SrtFile{filename: v}
		srtFiles = append(srtFiles, srt)
		go func(srt *SrtFile, c chan int) {
			err = srt.Read()
			if err != nil {
				fmt.Println("ReadStrFile:", err)
				os.Exit(1)
			}
			c <- 1
		}(srt, wait)
	}

	for i := 0; i < cap(wait); i++ {
		<-wait
	}

	mergedSrt, err := MergeSrt(srtFiles, timeOffset)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mergedSrt.filename = outputSub
	mergedSrt.Write()
}
