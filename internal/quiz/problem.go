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

func start(ctx context.Context, p []Problem, pc chan *Problem) {
	pc <- &p[0]
	i := 1
	for {
		select {
		case pc <- &p[i]:
			i++
			if i == len(p) {
				close(pc)
				return
			}
		case <-ctx.Done():
			fmt.Println("Problemer Closed")
			close(pc)
			return
		}
	}
}

func StartFromFile(ctx context.Context, filename string) (problemChan chan *Problem, err error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	p := parseLines(lines)
	pc := make(chan *Problem)
	go start(ctx, p, pc)

	return pc, nil
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
