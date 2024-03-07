package main

import (
	"os"
	"testing"
	"time"
)

func createTempProblemsFile(dir string, problems []problem) (*os.File, error) {
	file, err := os.CreateTemp(dir, "problems*.csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content := ""
	for i, p := range problems {
		content += p.Question
		content += ","
		content += p.Answer

		if i != len(problems) {
			content += "\n"
		}
	}

	file.WriteString(content)

	return file, nil
}

func TestReadProblems(t *testing.T) {
	expectedProblems := []problem{
		{"5+5", "10"},
		{"1+1", "2"},
	}
	file, err := createTempProblemsFile(t.TempDir(), expectedProblems)
	if err != nil {
		t.Fatal(err)
	}

	actualProblems, err := readProblems(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	for i, actual := range actualProblems {
		expected := expectedProblems[i]
		if expected.Question != actual.Question {
			t.Errorf("got %s, want %s", actual.Question, expected.Question)
		}
		if expected.Answer != actual.Answer {
			t.Errorf("got %s, want %s", actual.Answer, expected.Answer)
		}
	}
}

func TestAnswerProblem(t *testing.T) {
	expectedProblem := problem{"5+5", "10"}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(expectedProblem.Answer))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	defer func(f *os.File) { os.Stdin = f }(os.Stdin)
	os.Stdin = r
	defer r.Close()

	answerCh := make(chan string)
	go answerProblem(expectedProblem, 0, answerCh)

	actualAnswer := <-answerCh
	if actualAnswer != expectedProblem.Answer {
		t.Errorf("got %s, want %s", actualAnswer, expectedProblem.Answer)
	}
}

func TestAnswerProblems(t *testing.T) {
	expectedProblems := []problem{
		{"5+5", "10"},
		{"1+1", "2"},
	}
	ans := ""
	for _, p := range expectedProblems {
		ans += p.Answer
		ans += "\n"
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	_, err = w.Write([]byte(ans))
	if err != nil {
		t.Fatal(err)
	}
	w.Close()

	defer func(f *os.File) { os.Stdin = f }(os.Stdin)
	os.Stdin = r
	defer r.Close()

	correct := answerProblems(expectedProblems, 30*time.Second)
	if correct != len(expectedProblems) {
		t.Errorf("got %d, want %d", correct, len(expectedProblems))
	}
}
