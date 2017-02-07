package main

import (
	"container/list"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputSubsStr  string
	timeOffsetStr string

	inputSubs  []string
	outputSub  string
	timeOffset []string
)

func parseOptions() error {
	if inputSubsStr == "" || outputSub == "" || timeOffsetStr == "" {
		return errors.New("invalid parameters")
	}
	inputSubs = strings.Split(inputSubsStr, ";")
	timeOffset = strings.Split(timeOffsetStr, ";")

	fmt.Println("input sub files:", inputSubs)
	fmt.Println("output sub file:", outputSub)
	fmt.Println("time offset:", timeOffset)

	return nil
}

func main() {
	flag.StringVar(&inputSubsStr, "i", "", "input subtitle files, separated by ';'.")
	flag.StringVar(&outputSub, "o", "", "ouput subtitle file.")
	flag.StringVar(&timeOffsetStr, "t", "", "time offset, separated by ';'.")
	flag.Parse()

	err := parseOptions()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	srtFileNum := len(inputSubs)
	srtLists := make([]*list.List, srtFileNum)

	for i, v := range inputSubs {
		srtLists[i], err = ReadSrtFile(v)
		if err != nil {
			fmt.Println("ReadStrFile:", err)
			os.Exit(1)
		}
	}

	MergeSrt(srtLists, timeOffset)

	WriteSrtFile(srtLists, outputSub)
}
