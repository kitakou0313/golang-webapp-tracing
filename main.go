package main

import (
	"context"
	"go-tracing/cmd/app"
	"log"
	"os"
	"os/signal"
)

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
