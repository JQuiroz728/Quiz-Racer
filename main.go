package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

type problem struct {
	question string
	answer   string
}

func main() {
	filename := flag.String("csv", "quiz.csv", "csv file in 'question, answer' format")
	timeLimit := flag.Int("limit", 30, "time limit for quiz in seconds")
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		exit(fmt.Sprintf("Failed to open file: %s\n", *filename))
	}

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		exit("Failed to parse the provided file")
	}

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	problems := parseLines(lines)
	correct := 0
	// Problem loop
	for i, prob := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, prob.question)
		// Routine
		answerChannel := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChannel <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("\nTime Elapsed. You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerChannel:
			if answer == prob.answer {
				correct++
			}
		}
	}
	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []problem {
	retrn := make([]problem, len(lines)) // Return slice of problems. Length equal to total number of lines in .csv file. Assume every row in .csv file is a problem.
	for i, line := range lines {
		retrn[i] = problem{
			question: line[0],
			answer:   strings.TrimSpace(line[1]),
		}
	}
	return retrn
}

func exit(message string) {
	fmt.Println(message)
	os.Exit(1)
}
