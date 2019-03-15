package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/raonifn/gophercises/internal/quiz"
)

var filename string
var timeout time.Duration

func main() {
	flag.StringVar(&filename, "questions", "problems.csv", "Questions filename")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Timeout to answer all questions")
	flag.Parse()

	ctx, cancelFunc := context.WithCancel(context.Background())
	pc, nc, err := quiz.StartFromFile(filename, ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	anwsChan := quiz.CreateAsker(pc, ctx)

	var question quiz.Problem

	fmt.Printf("Correct answers: %d\n", correct)
}

func runQuestions(anwsChan chan string) (correct int) {
	timer := time.After(timeout)
	correct = 0
	for {
		select {
		case anwser, ok := <-anwsChan:
			if !ok {
				return correct
			}
			if anwser == question.A {
				correct++
			}
			nc <- true
		case <-timer:
			fmt.Println("Your time has over")
			cancelFunc()
			return correct
		}
	}
}
