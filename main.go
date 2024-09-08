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

func handleAny(context echo.Context) error {
	targetPostfixURL := context.Request().URL.String()
	targetURL := string(config.AppConfig.ProxyTargetUrl) + targetPostfixURL
	req, err := http.NewRequest(context.Request().Method, targetURL, context.Request().Body)
	if err != nil {
		return context.String(http.StatusBadRequest, "Bad Request")
	}
	req.Header = context.Request().Header
	println("REQUEST URL: %s", targetURL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return context.String(http.StatusBadGateway, "Failed to connect to destination server")
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		context.Response().Header().Set(k, v[0])
	}
	context.Response().WriteHeader(resp.StatusCode)

	_, err = io.Copy(context.Response().Writer, resp.Body)
	return err
}
