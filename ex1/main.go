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

// name of csv file in the directory
var csvfile string

// number of seconds to run the quiz
var limit int

var totalQuestions int

// Represents a question and answer
type Question struct {
	Question string
	Answer   string
}

func main() {
	flag.StringVar(&csvfile, "csv", "problems.csv", "A csv file in the format of 'question,answer'")
	flag.IntVar(&limit, "limit", 30, "Time limit for the quiz in seconds")

	flag.Parse()

	var questions []Question
	ParseQuestions(&questions)
	totalQuestions = len(questions)

	resultChan := make(chan int)

	// Run goroutine to listen for result
	go func() {
		result := 0
		for {
			select {
			case res := <-resultChan:
				result = res
			case <-time.After(time.Duration(limit) * time.Second):
				fmt.Printf("\nYou scored %d out of %d.\n", result, totalQuestions)
				os.Exit(1)
			}
		}
	}()

	index := 1
	totalCorrect := 0
	reader := bufio.NewReader(os.Stdin)
	// Go through each question
	for question := range questions {
		q := questions[question]
		fmt.Printf("Problem #%d: %s = ", index, q.Question)

		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == q.Answer {
			totalCorrect++
			resultChan <- totalCorrect
		}
		index++
	}

	//fmt.Println(questions)

}

func ParseQuestions(questions *[]Question) {
	// Open csv file
	file, err := os.Open(csvfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		*questions = append(*questions, Question{line[0], line[1]})
	}
}
