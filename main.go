package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

const problemsFilename = "problems.csv"

type Problem struct {
	Question string
	Answer   string
}

// readProblems reads all problems from a CSV file to a Problem slice.
func readProblems(filename string) ([]Problem, error) {
	problems := make([]Problem, 0)

	file, err := os.Open(filename)
	if err != nil {
		return problems, err
	}

	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return problems, err
	}

	for _, record := range records {
		question := record[0]
		answer := record[1]
		problems = append(problems, Problem{question, answer})
	}

	return problems, nil
}

// answerProblems prints each question of the problem and reads the user answer
// from standard input. It returns the number of problems answered correctly.
func answerProblems(problems []Problem) int {
	correct := 0

	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.Question)

		var answer string
		fmt.Scanln(&answer)

		if answer == p.Answer {
			correct++
		}
	}

	return correct
}

func main() {
	problems, err := readProblems(problemsFilename)
	if err != nil {
		log.Fatalf("read problems: %v", err)
	}

	correct := answerProblems(problems)
	fmt.Printf("You scored %d out of %d\n", correct, len(problems))
}
