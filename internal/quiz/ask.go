package quiz

import (
	"context"
	"fmt"
)

type asker struct {
	askChan   chan *Problem
	answChan  chan string
	questions int16
}

func (a *asker) start(ctx context.Context) {
	for {
		select {
		case p, ok := <-a.askChan:
			if !ok {
				close(a.answChan)
				return
			}

			var answer string
			a.questions++
			fmt.Printf("Problem: #%d: %s = ", a.questions, p.Q)

			fmt.Scanf("%s\n", &answer)
			a.answChan <- answer
		case <-ctx.Done():
			fmt.Println("Asker closed")
			close(a.answChan)
			return
		}
	}
}

func CreateAsker(askChan chan *Problem, ctx context.Context) chan string {
	answChan := make(chan string)

	a := asker{askChan, answChan, 0}

	go a.start(ctx)

	return answChan
}
