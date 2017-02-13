package main

import (
	"bufio"
	"bytes"
	"container/list"
	"strings"
	"testing"
)

const (
	srt1 = `1
00:00:00,000 --> 00:00:11,111
First Line
First Line2


2
00:00:22,222 --> 00:00:33,333
Second
Second 2

`
	srt2 = `1
00:00:01,001 --> 00:00:02,002
Third
Third 2

2
00:00:03,003 --> 00:00:04,004
Four

`

	// change Counter and TimeCode
	srt2_1 = `3
00:10:21,301 --> 00:10:22,302
Third
Third 2

4
00:10:23,303 --> 00:10:24,304
Four

`
	srt1_dos = "1\r\n00:00:00,000 --> 00:00:11,111\r\nFirst Line\r\nFirst Line2\r\n\r\n\r\n" +
		"2\r\n00:00:22,222 --> 00:00:33,333\r\nSecond\r\nSecond 2\r\n\r\n"

	srt2_dos = "1\r\n00:00:01,001 --> 00:00:02,002\r\nThird\r\nThird 2\r\n\r\n" +
		"2\r\n00:00:03,003 --> 00:00:04,004\r\nFour\r\n"

	srt2_1_dos = "3\r\n00:10:21,301 --> 00:10:22,302\r\nThird\r\nThird 2\r\n\r\n" +
		"4\r\n00:10:23,303 --> 00:10:24,304\r\nFour\r\n\r\n"
)

func Test_Str2TimeCode(t *testing.T) {
	str := "12:34:56,789"
	var tc TimeCode
	Str2TimeCode(str, &tc)

	newStr := tc.String()
	if newStr != str {
		t.Error("tc String wrong:", newStr)
	}
}

func Test_Add(t *testing.T) {
	str := "11:11:11,111"
	var tc TimeCode
	Str2TimeCode(str, &tc)

	tc.Add(str)
	newStr := tc.String()
	if newStr != "22:22:22,222" {
		t.Error("tc String wrong:", newStr)
	}
}

func Test_Merge(t *testing.T) {
	offset := []string{"00:10:20,300"}
	lineEnd = "\n"

	reader1 := bufio.NewReader(strings.NewReader(srt1))
	reader2 := bufio.NewReader(strings.NewReader(srt2))

	list1, err := doReadSrt(reader1)
	if err != nil {
		t.Fatal(err)
	}
	list2, err := doReadSrt(reader2)
	if err != nil {
		t.Fatal(err)
	}
	lists := []*list.List{list1, list2}

	doMergeSrt(lists, offset)

	var mergedSrt bytes.Buffer
	writer := bufio.NewWriter(&mergedSrt)
	err = doWriteSrt(lists, writer)
	if err != nil {
		t.Error(err)
	}
	mergedString := mergedSrt.String()
	if mergedString != srt1+srt2_1 {
		t.Error("merged failed")
	}
}

func Test_MergeDos(t *testing.T) {
	offset := []string{"00:10:20,300"}
	lineEnd = "\r\n"

	reader1 := bufio.NewReader(strings.NewReader(srt1_dos))
	reader2 := bufio.NewReader(strings.NewReader(srt2_dos))

	list1, err := doReadSrt(reader1)
	if err != nil {
		t.Fatal(err)
	}
	list2, err := doReadSrt(reader2)
	if err != nil {
		t.Fatal(err)
	}
	lists := []*list.List{list1, list2}

	doMergeSrt(lists, offset)

	var mergedSrt bytes.Buffer
	writer := bufio.NewWriter(&mergedSrt)
	err = doWriteSrt(lists, writer)
	if err != nil {
		t.Error(err)
	}
	mergedString := mergedSrt.String()
	if mergedString != srt1_dos+srt2_1_dos {
		t.Error("merged failed")
	}
}
