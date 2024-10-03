package httpHadnler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var loger = logrus.New()

// Возвращает response клиенту
func ResponseReturner(code int, resp string, err error, c *gin.Context) {
	if err != nil {
		loger.Warn(err)
		c.JSON(code, gin.H{"reason": err.Error()})
	} else {
		loger.Info(resp)
		c.String(code, resp)
	}
}

// Чтение параметров из тела запроса
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

// Чтение параметров запроса
func ReadQueryParams(c *gin.Context) map[string]string {
	paramsData := make(map[string]string)

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			paramsData[key] = values[0]
		}
	}

	return paramsData
}

// Перевод какого-либо объекта в JSON
func ToJSON(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		loger.Warn("JSON Unmarshal Error - ", err.Error())
		return "", errors.New("internal server error")
	}
	return string(jsonData), nil
}
