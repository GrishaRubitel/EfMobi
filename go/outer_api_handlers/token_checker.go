package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"time"
)

type activeToken struct {
	token        string
	clientId     string
	clientSecret string
	timer        int64
}

func (aT *activeToken) initToken() error {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "43f56af05fbd430a9e36aaad45cc5bc3")
	data.Set("client_secret", "beb497b4a876446d8cc8ace315413ecb")

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		loger.Fatal("error while creating Spotify API request - ", err.Error())
		return errors.New("internal server error")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	req.URL.RawQuery = data.Encode()

	resp, err := client.Do(req)
	if err != nil {
		loger.Fatal("error while calling spotify data - ", err.Error())
		return errors.New("internal server error")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		loger.Fatal("error while parsing access token response - ", err.Error())
		return errors.New("internal server error")
	}

	m := make(map[string]interface{})

	err = json.Unmarshal(body, &m)
	if err != nil {
		loger.Fatal("error decoding JSON - ", err.Error())
		return errors.New("internal server error")
	}

	aT.token = m["access_token"].(string)
	loger.Info("new spotify access token - " + aT.token)
	if t, ok := m["expires_in"]; ok {
		aT.timer = time.Now().Unix() + int64(t.(float64))
	} else {
		aT.timer = 0
	}
	return nil
}

func (aT *activeToken) checkToken() (string, error) {
	if aT.timer <= time.Now().Unix() {
		if err := aT.initToken(); err != nil {
			return "", err
		}
	}
	return aT.token, nil
}
