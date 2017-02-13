package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
 * The timecode format used is hours:minutes:seconds,milliseconds with time
 * units fixed to two zero-padded digits and fractions fixed to three
 * zero-padded digits (00:00:00,000). The fractional separator used is the comma.
 */
type TimeCode struct {
	time time.Time
}

func parseTimeCodeStr(str string) (hour, min, sec, msec int) {
	hour, _ = strconv.Atoi(str[0:2])
	min, _ = strconv.Atoi(str[3:5])
	sec, _ = strconv.Atoi(str[6:8])
	msec, _ = strconv.Atoi(str[9:12])
	return
}

func Str2TimeCode(str string, t *TimeCode) {
	hour, min, sec, msec := parseTimeCodeStr(str)
	nanosec := msec * 1000 * 1000
	t.time = time.Date(2000, time.April, 1, hour, min, sec, nanosec, time.UTC)
}

func (t TimeCode) String() string {
	str := t.time.Format("15:04:05.000")
	return str[0:8] + "," + str[9:]
}

func (t *TimeCode) Add(str string) error {
	hour, min, sec, msec := parseTimeCodeStr(str)
	dstr := fmt.Sprintf("%dh%dm%ds%dms", hour, min, sec, msec)
	d, err := time.ParseDuration(dstr)
	if err != nil {
		fmt.Println("parse duration:", err)
		return err
	}
	t.time = t.time.Add(d)
	return nil
}

type SrtEntry struct {
	Counter   int
	StartTime TimeCode
	EndTime   TimeCode
	Text      string
}

func (e SrtEntry) String() string {
	return fmt.Sprintf("%d%s%v --> %v%s%s%s", e.Counter, lineEnd, e.StartTime,
		e.EndTime, lineEnd, e.Text, lineEnd)
}

type SrtFile struct {
	filename  string
	entryList *list.List
}

type SrtMerged struct {
	filename       string
	entryListArray []*list.List
}

func readLine(r *bufio.Reader) (string, error) {
	str, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}

	str = strings.Trim(str, "\r\n")
	return str, nil
}

func processText(reader *bufio.Reader, entry *SrtEntry) (string, error) {
	for {
		str, err := readLine(reader)
		if err != nil {
			return "", err
		}

		if str != "" {
			entry.Text = entry.Text + str + lineEnd
		} else {
			for {
				nextLine, err := readLine(reader)
				if err != nil {
					return "", err
				}

				if nextLine != "" {
					return nextLine, nil
				}
				entry.Text = entry.Text + str + lineEnd
			}
		}
	}
}

func doReadSrt(reader *bufio.Reader) (*list.List, error) {
	l := list.New()
	counterLine, err := readLine(reader)
	if err != nil {
		return nil, err
	}

	for {
		entry := &SrtEntry{}

		// process Counter
		entry.Counter, err = strconv.Atoi(counterLine)
		if err != nil {
			fmt.Println("parse Counter:", err)
			break
		}

		// process TimeCode
		str, err := readLine(reader)
		if err != nil {
			break
		}
		Str2TimeCode(str[0:12], &entry.StartTime)
		Str2TimeCode(str[17:29], &entry.EndTime)

		// process Text
		counterLine, err = processText(reader, entry)

		if debug {
			fmt.Println("--------------------")
			fmt.Printf("%v", entry)
			fmt.Println("--------------------")
		}
		l.PushBack(entry)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
	}

	return l, nil
}

func (fi *SrtFile) Read() error {
	f, err := os.Open(fi.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	fi.entryList, err = doReadSrt(reader)
	if err != nil {
		return err
	}
	return nil
}

func doWriteSrt(entryListArray []*list.List, writer *bufio.Writer) error {
	for _, l := range entryListArray {
		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Fprint(writer, e.Value)
		}
	}
	writer.Flush()

	return nil
}

func (sm *SrtMerged) Write() error {
	f, err := os.Create(sm.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	return doWriteSrt(sm.entryListArray, writer)
}

func doMergeSrt(lists []*list.List, timeOffset []string) error {
	len := len(timeOffset)
	lastCounter := lists[0].Len()
	for i := 1; i <= len; i++ {
		for e := lists[i].Front(); e != nil; e = e.Next() {
			entry := e.Value.(*SrtEntry)
			entry.Counter = entry.Counter + lastCounter
			entry.StartTime.Add(timeOffset[i-1])
			entry.EndTime.Add(timeOffset[i-1])
		}
	}

	return nil
}

func MergeSrt(srtArray []*SrtFile, timeOffset []string) (*SrtMerged, error) {
	mergedSrt := &SrtMerged{}
	for _, v := range srtArray {
		mergedSrt.entryListArray = append(mergedSrt.entryListArray, v.entryList)
	}
	doMergeSrt(mergedSrt.entryListArray, timeOffset)
	return mergedSrt, nil
}
