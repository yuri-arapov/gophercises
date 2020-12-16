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

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	console := bufio.NewReader(os.Stdin)

	type Q struct {
		q, a string
		done bool
	}

	questionCh := make(chan Q)
	answerCh := make(chan string)

	csv := csv.NewReader(f)

	total, correct, line := 0, 0, 0

loop:
	for {
		line++

		go func() {
			data, err := csv.Read()
			if err == io.EOF {
				questionCh <- Q{done: true}
				return
			}
			if err != nil {
				log.Fatal(err)
			}
			if len(data) != 2 {
				log.Fatal(fmt.Errorf("invalid CSV line %d: 2 fields expected", line))
			}
			questionCh <- Q{q: data[0], a: data[1]}
		}()

		var q Q

		select {
		case <-timer.C:
			break loop // timeout, exit the loop
		case q = <-questionCh:
		}

		if q.done {
			break loop // all done, exit the loop
		}

		fmt.Printf("%s = ", q.q)
		total++

		go func() {
			input, err := console.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			answerCh <- input
		}()

		select {
		case <-timer.C:
			fmt.Println()
			break loop
		case input := <-answerCh:
			if strings.TrimSpace(input) == q.a {
				correct++
			}
		}
	}

	fmt.Printf("%d correct answers out of %d\n", correct, total)
}
