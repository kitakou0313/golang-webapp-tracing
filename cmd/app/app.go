package app

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, err
	}

	nStr := strconv.FormatUint(uint64(n), 10)
	span.SetAttributes(attribute.String(
		"repeat.n", nStr,
	))

	return n, err
}

func (a *App) WriteNthFibNum(ctx context.Context, n uint) {
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	ctx, span := otel.Tracer(name).Start(ctx, "WriteNthFibNum")
	defer span.End()
	span.SetAttributes(attribute.String("currentFunc", "Write"))

	span.AddEvent("Wating random seconds...")
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	span.AddEvent("Restart!")

	f, err := func(ctx context.Context) (uint64, error) {
		currentCtx, span := otel.Tracer(name).Start(ctx, "Fib")
		defer span.End()

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		fib, err := Fib(n)
		if err != nil {
			currentSpan := trace.SpanFromContext(currentCtx)
			currentSpan.RecordError(err)
			currentSpan.SetStatus(codes.Error, err.Error())
		}
		return fib, err
	}(ctx)

	if err != nil {
		a.l.Printf("Fibonacci(%d): %v\n", n, err)
	}
	a.l.Printf("Fibonacci(%d) = %d\n", n, f)

}
