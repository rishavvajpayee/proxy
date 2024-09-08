package main

import (
	"io"
	"net/http"
	config "proxyserver/pkg"
	"time"

	"github.com/labstack/echo/v4"
)

var httpClient *http.Client

func main() {
	config.LoadConfig()
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	e := echo.New()
	e.Any("/*", handleAny)
	e.Logger.Fatal(e.Start(":8080"))
}

func handleAny(c echo.Context) error {
	targetPostfixURL := c.Request().URL.String()
	targetURL := config.AppConfig.ProxyTargetUrl + targetPostfixURL
	req, err := http.NewRequest(c.Request().Method, targetURL, c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request: "+err.Error())
	}
	req.Header = c.Request().Header.Clone()
	resp, err := httpClient.Do(req)
	if err != nil {
		return c.String(http.StatusBadGateway, "Failed to connect to destination server: "+err.Error())
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		for _, value := range v {
			c.Response().Header().Add(k, value)
		}
	}
	c.Response().WriteHeader(resp.StatusCode)
	_, err = io.Copy(c.Response().Writer, resp.Body)
	return err
}
