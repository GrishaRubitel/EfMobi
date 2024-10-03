package main

import (
	"bytes"
	httpH "choomandco/efimobi/http_handler"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var loger = logrus.New()

func main() {
	envMap, err := godotenv.Read("../../.env")
	if err != nil {
		loger.Fatal("Error loading .env file - " + err.Error())
	}

	router := gin.Default()
	router.Use(cors.Default())
	aT := activeToken{
		clientId:     envMap["SPOTIFY_CLIENT_ID"],
		clientSecret: envMap["SPOTIFY_CLIENT_SECRET"],
		timer:        0,
	}

	router.GET("/token", func(c *gin.Context) {
		if token, err := aT.checkToken(); err != nil {
			httpH.ResponseReturner(http.StatusInternalServerError, "", err, c)
		} else {
			loger.Info("access token - " + token)
			httpH.ResponseReturner(http.StatusOK, token, nil, c)
		}
	})

	router.GET("/search_track", func(c *gin.Context) {
		title, _ := url.QueryUnescape(c.Query("title"))
		artist, _ := url.QueryUnescape(c.Query("artist"))
		if title == "" || artist == "" {
			er := "both title and artist parameters are required"
			loger.Warn(er)
			httpH.ResponseReturner(http.StatusBadRequest, "", errors.New(er), c)
		} else {
			if code, resp, err := searchTrack(title, artist, aT, envMap["FREE_PROXY_ADDRESS"]); err == nil {
				if resp, err := PrettyPrintJSON(resp); err == nil {
					httpH.ResponseReturner(code, resp, nil, c)
				} else {
					httpH.ResponseReturner(http.StatusInternalServerError, "", err, c)
				}
			} else {
				httpH.ResponseReturner(code, "", err, c)
			}

		}
	})

	router.Run(envMap["TOKEN_FINDER_ADDRESS"])
}

func PrettyPrintJSON(rawJSON []byte) (string, error) {
	var prettyJSON bytes.Buffer

	err := json.Indent(&prettyJSON, rawJSON, "", "    ")
	if err != nil {
		return "", fmt.Errorf("error formatting JSON: %w", err)
	}

	return prettyJSON.String(), nil
}
