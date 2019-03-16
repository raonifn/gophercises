package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/raonifn/gophercises/internal/quiz"
)

var (
	filename string
	timeout  time.Duration
	shuf     bool
)

func main() {
	flag.StringVar(&filename, "questions", "problems.csv", "Questions filename")
	flag.DurationVar(&timeout, "timeout", 30*time.Second, "Timeout to answer all questions")
	flag.BoolVar(&shuf, "shuffle", false, "Shuffle questions")
	flag.Parse()

	ctx, cancelFunc := context.WithCancel(context.Background())
	pc, err := quiz.StartFromFile(ctx, filename, shuf)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	intro()

	asker := quiz.CreateAsker(ctx, pc)
	correct := runQuestions(asker, cancelFunc)

	fmt.Printf("Correct answers: %d\n", correct)
}

func intro() {
	fmt.Printf("Press enter to start. Timeout %v", timeout)
	fmt.Scanln()
}

func runQuestions(asker quiz.Asker, cancelFunc context.CancelFunc) (correct int) {
	timer := time.After(timeout)
	correct = 0
	for {
		select {
		case answer, ok := <-asker.AnswChan:
			if !ok {
				return correct
			}
			if answer.IsCorrect() {
				correct++
			}
		case <-timer:
			fmt.Println("Your time has over")
			cancelFunc()
			return correct
		}
	}
}
