package main

import (
	"io"
	"net/http"
	config "proxyserver/pkg"

	"github.com/labstack/echo/v4"
)

func main() {
	config.LoadConfig()
	e := echo.New()
	e.Any("/*", handleAny)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleAny(c echo.Context) error {
	targetPostfixURL := c.Request().URL.String()
	targetURL := string(config.AppConfig.ProxyTargetUrl) + targetPostfixURL

	req, err := http.NewRequest(c.Request().Method, targetURL, c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}
	req.Header = c.Request().Header

	responseChannel := make(chan *http.Response)
	errorChannel := make(chan error)

	go func() {
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			errorChannel <- err
			return
		}
		responseChannel <- resp
	}()

	select {
	case resp := <-responseChannel:
		defer resp.Body.Close()

		for k, v := range resp.Header {
			c.Response().Header().Set(k, v[0])
		}
		c.Response().WriteHeader(resp.StatusCode)

		_, err := io.Copy(c.Response().Writer, resp.Body)
		return err

	case err := <-errorChannel:
		return c.String(http.StatusBadGateway, "Failed to connect to destination server: "+err.Error())
	}
}
