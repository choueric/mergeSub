package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strings"
	"testing"
)

func Test_Str2TimeCode(t *testing.T) {
	str := "12:34:56,789"
	var tc TimeCode
	Str2TimeCode(str, &tc)
	fmt.Println(tc)
	fmt.Println(&tc)
}

func Test_Add(t *testing.T) {
	str := "11:11:11,111"
	var tc TimeCode
	Str2TimeCode(str, &tc)
	fmt.Println(tc)

	tc.Add(str)
	fmt.Println(tc)
}

func Test_Merge(t *testing.T) {
	srt1 := "1\n00:00:00,000 --> 00:00:11,111\nFirst Line\n\n2\n00:00:22,222 --> 00:00:33,333\nSecond\n\n"
	srt2 := "1\n00:00:01,001 --> 00:00:02,002\nThird\n\n2\n00:00:03,003 --> 00:00:04,004\nFour\n\n"
	offset := []string{"00:10:20,300"}

	fmt.Printf("----------- srt1 -------------\n")
	fmt.Print(srt1)
	fmt.Printf("----------- srt2 -------------\n")
	fmt.Print(srt2)

	reader1 := bufio.NewReader(strings.NewReader(srt1))
	reader2 := bufio.NewReader(strings.NewReader(srt2))

	list1, err := doReadSrt(reader1)
	if err != nil {
		t.Error(err)
	}
	list2, err := doReadSrt(reader2)
	if err != nil {
		t.Error(err)
	}

	lists := []*list.List{list1, list2}

	MergeSrt(lists, offset)

	fmt.Println("Offset:", offset)
	fmt.Printf("------------ merged ----------\n")
	writer := bufio.NewWriter(os.Stdout)
	err = doWriteSrt(lists, writer)
	if err != nil {
		t.Error(err)
	}

}
