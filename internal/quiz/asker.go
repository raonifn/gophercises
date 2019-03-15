package quiz

import (
	"context"
	"fmt"
)

type Asker struct {
	askChan   chan *Problem
	AnswChan  chan *AskedQuestion
	Questions int16
}

type AskedQuestion struct {
	Problem *Problem
	Answer  string
}

func (a *Asker) start(ctx context.Context) {
	for {
		select {
		case p, ok := <-a.askChan:
			if !ok {
				close(a.AnswChan)
				return
			}

			var answer string
			a.Questions++
			fmt.Printf("Problem: #%d: %s = ", a.Questions, p.Q)

			fmt.Scanf("%s\n", &answer)
			a.AnswChan <- &AskedQuestion{Problem: p, Answer: answer}
		case <-ctx.Done():
			fmt.Println("Asker closed")
			close(a.AnswChan)
			return
		}
	}
}

func CreateAsker(ctx context.Context, askChan chan *Problem) Asker {
	AnswChan := make(chan *AskedQuestion)

	a := Asker{askChan, AnswChan, 0}

	go a.start(ctx)

	return a
}
