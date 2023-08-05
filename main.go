package main

import (
	"context"
	"go-tracing/cmd/app"
	"io"
	"log"
	"os"
	"os/signal"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		stdouttrace.WithPrettyPrint(),
	)

}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("fib"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)

	return r
}

func main() {
	// e := echo.New()

	// e.GET("/service-a-endpoint", func(c echo.Context) error {
	// 	url := "http://service-b:8080/service-b-endpoint"

	// 	for i := 0; i < 10; i++ {
	// 		resp, err := http.Get(url)
	// 		if err != nil {
	// 			e.Logger.Error(err.Error())
	// 			return c.String(http.StatusInternalServerError, "Error:"+err.Error())
	// 		}
	// 		defer resp.Body.Close()

	// 		byteArray, _ := ioutil.ReadAll(resp.Body)
	// 		e.Logger.Info((string(byteArray))) // htmlをstringで取得
	// 	}

	// 	return c.String(http.StatusOK, "Hello from service-a!")
	// })
	// e.GET("/service-b-endpoint", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello from service-b!")
	// })

	// e.Logger.Fatal(e.Start(":8080"))

	l := log.New(os.Stdout, "", 0)

	f, err := os.Create("trace.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()

	exp, err := newExporter(f)
	if err != nil {
		l.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			l.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	errCh := make(chan error)

	app := app.NewApp(os.Stdin, l)

	go func() {
		errCh <- app.Run(context.Background())
	}()

	select {
	case <-sigCh:
		l.Println("\nGoodBye")
		return
	case err := <-errCh:
		if err != nil {
			l.Fatal(err)
		}
	}
}
