package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	fname := "problems.csv"

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	//	scanner := bufio.NewScanner(f)
	//	for scanner.Scan() {
	//		fmt.Println(scanner.Text())
	//	}

	console := bufio.NewReader(os.Stdin)

	csv := csv.NewReader(f)
	total, correct, line := 0, 0, 0
	for {
		line++
		data, err := csv.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(data) != 2 {
			log.Fatal(fmt.Errorf("invalid CSV line %d: 2 fields expected", line))
		}
		q := data[0]
		a := data[1]
		fmt.Printf("%s = ", q)
		input, err := console.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		total++
		if strings.TrimSpace(input) == a {
			correct++
		}
	}
	fmt.Printf("%d correct answers out of %d\n", correct, total)
}
