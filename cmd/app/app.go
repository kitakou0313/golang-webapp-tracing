package app

import (
	"context"
	"fmt"
	"io"
	"log"
)

type App struct {
	r io.Reader
	l *log.Logger
}

func NewApp(r io.Reader, l *log.Logger) *App {
	return &App{
		r: r,
		l: l,
	}
}

func (a *App) Run(ctx context.Context) error {
	for {
		n, err := a.Poll(ctx)
		if err != nil {
			return err
		}

		a.WriteNthFibNum(ctx, n)
	}
}

func (a *App) Poll(ctx context.Context) (uint, error) {
	a.l.Print(
		"What Fib Number would you like to know:",
	)

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	return n, err
}

func (a *App) WriteNthFibNum(ctx context.Context, n uint) {
	f, err := Fib(n)
	if err != nil {
		a.l.Printf(
			"FibNum %d: %d\n", n, f,
		)
	} else {
		a.l.Printf(
			"FibNum %d: %v\n", n, err,
		)
	}
}
