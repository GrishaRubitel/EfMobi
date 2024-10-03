package httpHadnler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var loger = logrus.New()

func ResponseReturner(code int, resp string, err error, c *gin.Context) {
	if err != nil {
		loger.Warn(err)
		c.JSON(code, gin.H{"reason": err.Error()})
	} else {
		loger.Info(resp)
		c.String(code, resp)
	}
}

func ReadBodyData(c *gin.Context) (int, map[string]string, error) {
	var bodyData map[string]interface{}

	body, err := c.GetRawData()
	if err != nil {
		loger.Warn("Body processing error mk1 - " + err.Error())
		return http.StatusInternalServerError, nil, errors.New("error while processing client data")
	}

	if err := json.Unmarshal(body, &bodyData); err != nil {
		loger.Warn("Body processing error mk2 - " + err.Error())
		return http.StatusInternalServerError, nil, errors.New("error while processing client data")
	}

	result := make(map[string]string)

	for key, value := range bodyData {
		switch v := value.(type) {
		case string:
			result[key] = v
		case float64:
			result[key] = strconv.FormatFloat(v, 'f', -1, 64)
		default:
			loger.Warn("unsupported value type - ", key)
			return http.StatusBadRequest, nil, errors.New("unsupported value type - " + key)
		}
	}

	return http.StatusOK, result, nil
}

func ReadQueryParams(c *gin.Context) map[string]string {
	paramsData := make(map[string]string)

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			paramsData[key] = values[0]
		}
	}

	return paramsData
}

func ToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		loger.Warn("JSON Unmarshal Error - ", err.Error())
		return "", errors.New("internal server error")
	}
	return string(jsonData), nil
}

func GetNestedMap(m map[string]interface{}, keys ...string) (map[string]interface{}, error) {
	current := m
	for _, key := range keys {
		if val, ok := current[key]; ok {
			if nestedMap, ok := val.(map[string]interface{}); ok {
				current = nestedMap
			} else {
				return nil, fmt.Errorf("key '%s' does not contain a map", key)
			}
		} else {
			return nil, fmt.Errorf("key '%s' not found", key)
		}
	}
	return current, nil
}

func CallTokenAPI(url string) (int, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, "", fmt.Errorf("failed to make GET request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read response body: %w", err)
	}

	return resp.StatusCode, string(body), nil
}

func CallOuterApi(url string) (int, map[string]interface{}, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		loger.Panic("Error while creating API request - " + err.Error())
		return http.StatusInternalServerError, nil, errors.New("critical server error")
	}

	resp, err := client.Do(req)
	if err != nil {
		loger.Panic("Error while calling API - " + err.Error())
		return http.StatusInternalServerError, nil, errors.New("critical server error")
	}
	defer resp.Body.Close()

	code := resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, nil, errors.New("internal server error")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		loger.Warn("Error parsing JSON - ", err)
		return http.StatusInternalServerError, nil, errors.New("critical server error while parsing JSON")
	}

	return code, result, nil
}
