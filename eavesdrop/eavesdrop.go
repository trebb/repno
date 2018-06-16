package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"log"
	"os"
	"time"
)

var (
	helpMsg = "eavesdrop is a reverse-engineering tool that logs the communication between\n" +
		"touchscreen and mainboard of a Kawai digital piano.\n\n" +
		"You can type in short descriptions of the actions you're going to perform.\n" +
		"The descriptive line of text and the communication resulting from your actions\n" +
		"will be written to two files named like yyyymmddThhmmss_ui.txt\n" +
		"and yyyymmddThhmmss_mb.txt in the current directory.\n\n"
	uiPort = flag.String("ui", "/dev/ttyUSB0", "the serial interface connected to the touchscreen")
	mbPort = flag.String("mb", "/dev/ttyUSB1", "the serial interface connected to the mainboard")
)

const (
	tStamp         = "060102T15:04:05.000000"
	filenameTStamp = "20060102T150405"
)

func logRS232(port string, outFilename string, src string, hl chan string, msgs chan string) {
	s, err := serial.OpenPort(&serial.Config{Name: port, Baud: 115200})
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(outFilename + "_" + src + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf := make([]byte, 1)
	newLine := true
	bytesSent := 0
	for {
		t0 := time.Now()
		_, err := s.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Now()
		msgs <- string(src[0])
		bytesSent++
		if t.Sub(t0) > 100*time.Millisecond {
			msgs <- fmt.Sprintf("[%s:%d]\n", src, bytesSent)
			bytesSent = 0
			newLine = true
		}
		if newLine {
			sendHeadline(f, t, src, hl)
			fmt.Fprintf(f, "\n%-22v%s >>", t.Format(tStamp), src)
			newLine = false
		}
		fmt.Fprintf(f, "%3X", buf[0])
	}
}

func sendHeadline(f *os.File, t time.Time, src string, headlines chan string) {
	// write headline(s) if any
	for done := false; done == false; {
		select {
		case h := <-headlines:
			fmt.Fprintf(f, "\n%-22v%s #### %s ", t.Format(tStamp), src, h)
		default:
			done = true
		}
	}
}

func feedback(msgs chan string) {
	for {
		msg := <-msgs
		fmt.Print(msg)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, helpMsg)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	uiHeadlines := make(chan string, 100)
	mbHeadlines := make(chan string, 100)
	msg := make(chan string, 1000)
	outFilename := time.Now().Format(filenameTStamp)
	go logRS232(*uiPort, outFilename, "ui", uiHeadlines, msg)
	go logRS232(*mbPort, outFilename, "mb", mbHeadlines, msg)
	go feedback(msg)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		uiHeadlines <- scanner.Text()
		mbHeadlines <- scanner.Text()
	}
	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}
