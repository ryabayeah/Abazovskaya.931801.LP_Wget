package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WriteCounter struct {
	Total uint64
	Seconds int
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)

	return n, nil
}

func countSeconds(wc *WriteCounter){
	for {
		time.Sleep(time.Second)
		wc.Seconds++
		fmt.Println("Time:",wc.Seconds,"sec","|","Downloaded:", wc.Total,"byte")
	}
}

func main() {
	var fileUrl string
	fmt.Print("Link:")
	fmt.Scan(&fileUrl)
	fileName := fileUrl[1+strings.LastIndex(fileUrl, "/"):]
	err := DownloadFile(fileName, fileUrl)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error: no URL specified\n")
		return
	}
}

func DownloadFile(filepath string, url string) error {


	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		out.Close()
		return err
	}
	defer resp.Body.Close()
	fmt.Println("\nDownload started\nFile:",filepath,"\nTotal size:",resp.ContentLength)
	counter := &WriteCounter{}
	go countSeconds(counter)
	if _, err = io.Copy(out, io.TeeReader(resp.Body, counter)); err != nil {
		out.Close()
		return err
	}

	fmt.Print("\n")
	out.Close()

	if err = os.Rename(filepath+".tmp", filepath); err != nil {
		return err
	}
	fmt.Println("\nDownload finished\nFile:",filepath,"\nSize:",counter.Total,"byte","\nTime:",counter.Seconds,"sec\n")
	return nil
}