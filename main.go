package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	defaultProblemsFilename = "problems.csv"
	defaultTimeLimit        = 30
)

type problem struct {
	Question string
	Answer   string
}

// readProblems reads all problems from a CSV file to a problem slice.
func readProblems(filename string) ([]problem, error) {
	file, err := os.Open(filename)
	if err != nil {
		return make([]problem, 0), err
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return make([]problem, 0), err
	}

	problems := make([]problem, len(records))
	for i, record := range records {
		question := record[0]
		answer := strings.TrimSpace(record[1])
		problems[i] = problem{question, answer}
	}

	return problems, nil
}

// answerProblem prints the question of problem p at index i and reads the user answer
// from standard input. It sends the result of the user answer as a boolean to
// channel c.
func answerProblem(p problem, i int, c chan<- string) {
	fmt.Printf("Problem #%d: %s = ", i+1, p.Question)

	var answer string
	fmt.Scanln(&answer)

	c <- answer
}

// answerProblems prints each question of the problem and reads the user answer
// from standard input. It times out after limit seconds. It returns
// the number of problems answered correctly.
func answerProblems(problems []problem, limit time.Duration) int {
	correct := 0
	answerCh := make(chan string)
	timeout := time.After(limit)

	for i, p := range problems {
		go answerProblem(p, i, answerCh)

		select {
		case <-timeout:
			fmt.Println()
			return correct
		case answer := <-answerCh:
			if answer == p.Answer {
				correct++
			}
		}
	}

	return correct
}

func main() {
	filename := flag.String("csv", defaultProblemsFilename, "a csv file in the format of 'question,answer'")
	timeLimit := flag.Int("limit", defaultTimeLimit, "the time limit for the quiz in seconds")
	flag.Parse()

	problems, err := readProblems(*filename)
	if err != nil {
		log.Fatalf("read problems: %v", err)
	}

	correct := answerProblems(problems, time.Duration(*timeLimit)*time.Second)
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}
