package app

import (
	"context"
	"fmt"
	"io"
	"log"

	"go.opentelemetry.io/otel"
)

const name = "fib"

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
		newCtx, span := otel.Tracer(name).Start(ctx, "Run")

		n, err := a.Poll(newCtx)
		if err != nil {
			span.End()
			return err
		}

		a.WriteNthFibNum(newCtx, n)
		span.End()
	}
}

func (a *App) Poll(ctx context.Context) (uint, error) {
	_, span := otel.Tracer(name).Start(ctx, "Poll")
	defer span.End()

	a.l.Print(
		"What Fib Number would you like to know:",
	)

	var n uint
	_, err := fmt.Fscanf(a.r, "%d\n", &n)
	return n, err
}

func (a *App) WriteNthFibNum(ctx context.Context, n uint) {
	ctx, span := otel.Tracer(name).Start(ctx, "WriteNthFibNum")
	defer span.End()

	f, err := func(ctx context.Context) (uint64, error) {
		_, span := otel.Tracer(name).Start(ctx, "Fib")
		defer span.End()

		return Fib(n)
	}(ctx)

	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	}
	a.l.Printf("Fibonacci(%d) = %d\n", n, f)

}
