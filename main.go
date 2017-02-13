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

	fileNum := len(inputSubs)
	srtLists := make([]*list.List, fileNum)
	wait := make(chan int, fileNum)

	for i, v := range inputSubs {
		go func(index int, filename string, c chan int) {
			srtLists[index], err = ReadSrtFile(filename)
			if err != nil {
				fmt.Println("ReadStrFile:", err)
				os.Exit(1)
			}
			c <- 1
		}(i, v, wait)
	}

	for i := 0; i < fileNum; i++ {
		<-wait
	}

	MergeSrt(srtLists, timeOffset)

	WriteSrtFile(srtLists, outputSub)
}
