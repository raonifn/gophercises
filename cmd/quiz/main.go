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
	pc, err := quiz.StartFromFile(ctx, filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	asker := quiz.CreateAsker(ctx, pc)
	correct := runQuestions(asker, cancelFunc)

	fmt.Printf("Correct answers: %d\n", correct)
}

func runQuestions(asker quiz.Asker, cancelFunc context.CancelFunc) (correct int) {
	timer := time.After(timeout)
	correct = 0
	for {
		select {
		case anwser, ok := <-asker.AnswChan:
			if !ok {
				return correct
			}
			if anwser.Answer == anwser.Problem.A {
				correct++
			}
		case <-timer:
			fmt.Println("Your time has over")
			cancelFunc()
			return correct
		}
	}
}
