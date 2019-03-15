package quiz

import (
	"context"
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Problem struct {
	Q string
	A string
}

func normalizeString(str string) string {
	return strings.TrimSpace(strings.ToUpper(str))
}

func (p *Problem) IsCorrect(answer string) bool {
	return normalizeString(answer) == normalizeString(p.A)
}

func shuffle(p []Problem) []Problem {
	ret := make([]Problem, 0)
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(p)) {
		val := p[i]
		ret = append(ret, val)
	}
	return ret
}

func start(ctx context.Context, p []Problem, pc chan *Problem, shuf bool) {
	if shuf {
		p = shuffle(p)
	}
	i := 0
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

func StartFromFile(ctx context.Context, filename string, shuf bool) (problemChan chan *Problem, err error) {
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
	go start(ctx, p, pc, shuf)

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
