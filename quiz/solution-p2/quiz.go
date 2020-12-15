package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	var fname string
	flag.StringVar(&fname, "quiz-file", "problems.csv", "quiz data in CSV format")
	var timeLimit int
	flag.IntVar(&timeLimit, "time-limit", 30, "quiz time limit, seconds")
	flag.Parse()

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//	scanner := bufio.NewScanner(f)
	//	for scanner.Scan() {
	//		fmt.Println(scanner.Text())
	//	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	_ = timer

	console := bufio.NewReader(os.Stdin)

	eofCh := make(chan struct{})
	questionCh := make(chan []string)
	answerCh := make(chan string)

	csv := csv.NewReader(f)

	total, correct, line := 0, 0, 0

loop:
	for {
		line++

		go func() {
			data, err := csv.Read()
			if err == io.EOF {
				eofCh <- struct{}{}
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			questionCh <- data
		}()

		var q, a string

		select {
		case <-eofCh:
			break loop
		case data := <-questionCh:
			if len(data) != 2 {
				log.Fatal(fmt.Errorf("invalid CSV line %d: 2 fields expected", line))
			}
			q = data[0]
			a = data[1]
		case <-timer.C:
			break loop
		}

		fmt.Printf("%s = ", q)
		total++

		go func() {
			input, err := console.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			answerCh <- input
		}()

		select {
		case input := <-answerCh:
			if strings.TrimSpace(input) == a {
				correct++
			}
		case <-timer.C:
			fmt.Println()
			break loop
		}
	}

	fmt.Printf("%d correct answers out of %d\n", correct, total)
}
