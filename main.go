package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func newExporter(ctx context.Context) (trace.SpanExporter, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint("jaeger:4318"),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(ctx, client)

	return exporter, err
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("trace-with-http"),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)

	return r
}

func traceWithEcho() {
	e := echo.New()

	e.GET("/service-a-endpoint", func(c echo.Context) error {
		url := "http://service-b:8080/service-b-endpoint"

		for i := 0; i < 10; i++ {
			// resp, err := http.Get(url)
			resp, err := otelhttp.Get(c.Request().Context(), url)

			if err != nil {
				e.Logger.Error(err.Error())
				return c.String(http.StatusInternalServerError, "Error:"+err.Error())
			}
			defer resp.Body.Close()

			byteArray, _ := ioutil.ReadAll(resp.Body)
			e.Logger.Info((string(byteArray))) // htmlをstringで取得
		}

		return c.String(http.StatusOK, "Hello from service-a!")
	})
	e.GET("/service-b-endpoint", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from service-b!")
	})

	e.Use(otelecho.Middleware("instrumented-echo"))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status} header=${header:traceparent}\n",
	}))

	e.Logger.Fatal(e.Start(":8080"))

}

func fibWithTrace() {
	// l := log.New(os.Stdout, "", 0)

	// ctx := context.Background()

	// exp, err := newExporter(ctx)
	// if err != nil {
	// 	l.Fatal(err)
	// }

	// tp := trace.NewTracerProvider(
	// 	trace.WithBatcher(exp),
	// 	trace.WithResource(newResource()),
	// )
	// defer func() {
	// 	if err := tp.Shutdown(ctx); err != nil {
	// 		l.Fatal(err)
	// 	}
	// }()
	// otel.SetTracerProvider(tp)

	// sigCh := make(chan os.Signal, 1)
	// signal.Notify(sigCh, os.Interrupt)

	// errCh := make(chan error)

	// app := app.NewApp(os.Stdin, l)

	// go func() {
	// 	errCh <- app.Run(context.Background())
	// }()

	// select {
	// case <-sigCh:
	// 	l.Println("\nGoodBye")
	// 	return
	// case err := <-errCh:
	// 	if err != nil {
	// 		l.Fatal(err)
	// 	}
	// }
}

var tracer = otel.Tracer(
	"test-instrumented-libs",
)

func sleepy(ctx context.Context) {
	_, span := tracer.Start(ctx, "sleep")
	defer span.End()

	time.Sleep(1 * time.Second)
	span.SetAttributes(
		attribute.Int("sleep.duration", int(1*time.Second)),
	)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, I am instrumented automatically!")
	ctx := r.Context()
	sleepy(ctx)
}

func traceWithInstrumentedLibs() {
	handler := http.HandlerFunc(httpHandler)
	wrappedHandler := otelhttp.NewHandler(
		handler, "hello-instrumented",
	)
	http.Handle("/hello-instrumented", wrappedHandler)

	log.Fatal(http.ListenAndServe(
		":3030", nil,
	))
}

func main() {
	ctx := context.Background()

	exp, err := newExporter(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)

	traceWithEcho()
}
