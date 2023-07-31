package main

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/service-a-endpoint", func(c echo.Context) error {
		url := "http://service-b:8080/service-b-endpoint"

		for i := 0; i < 10; i++ {
			resp, err := http.Get(url)
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

	e.Logger.Fatal(e.Start(":8080"))
}
