package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Запрос трека в Spotify Web API
func searchTrack(title string, artist string, aT activeToken, proxyURL string) (int, []byte, error) {
	apiURL := fmt.Sprintf("https://api.spotify.com/v1/search?q=track%%3D%s%%26artist%%3D%s&type=track&limit=1&offset=0", url.QueryEscape(title), url.QueryEscape(artist))

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating request -", err)
		return http.StatusInternalServerError, nil, errors.New("internal server error")
	}

	token, err := aT.checkToken()
	if err != nil {
		log.Println(err.Error())
		return http.StatusInternalServerError, nil, err
	}
	req.Header.Add("Authorization", "Bearer "+token)

	var client *http.Client
	if proxyURL != "" {
		proxy, err := url.Parse(proxyURL)
		if err != nil {
			loger.Warn("invalid proxy URL -", err)
			return http.StatusBadGateway, nil, errors.New("invalid proxy URL")
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxy),
		}
		client = &http.Client{Transport: transport}
	} else {
		client = &http.Client{}
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error making request -", err)
		return http.StatusInternalServerError, nil, errors.New("internal server error")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response -", err)
		return http.StatusInternalServerError, nil, errors.New("internal server error")
	}

	return http.StatusOK, body, nil
}
