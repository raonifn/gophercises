package quiz

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

type Problem struct {
	Q string
	A string
}

func start(p []Problem, pc chan *Problem, nc chan bool, ctx context.Context) {
	pc <- &p[0]
	for {
		i := 1
		select {
		case <-nc:
			if i == len(p) {
				close(pc)
				close(nc)
				return
			}
			pc <- &p[i]
			i++
		case <-ctx.Done():
			fmt.Println("Problemer Closed")
			close(pc)
			close(nc)
			return
		}
	}
}

func StartFromFile(filename string, ctx context.Context) (problemChan chan *Problem, nextChan chan bool, err error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, nil, err
	}
	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		return nil, nil, err
	}

	p := parseLines(lines)
	pc := make(chan *Problem)
	nc := make(chan bool)
	go start(p, pc, nc, ctx)

	return pc, nc, nil
}

func parseLines(lines [][]string) []Problem {
	problems := []Problem{}
	for _, line := range lines {
		p := Problem{
			Q: strings.TrimSpace(line[0]),
			A: strings.TrimSpace(line[1]),
		}
		problems = append(problems, p)
	}
	return problems
}
