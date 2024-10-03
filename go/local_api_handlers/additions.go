package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

func extractStrings(lyrics string, limit int, offset int) string {
	var result []string
	var currentSubstring strings.Builder

	skipCount := 0
	capCount := 0
	first := true

	for _, char := range lyrics {
		if unicode.IsUpper(char) && !first {
			if skipCount < offset {
				skipCount++
			} else {
				result = append(result, currentSubstring.String())
				capCount++
				if limit > 0 && capCount >= limit {
					break
				}
			}
			currentSubstring.Reset()
		}

		currentSubstring.WriteRune(char)

		if first && unicode.IsUpper(char) {
			first = false
		}
	}

	if currentSubstring.Len() > 0 && skipCount >= offset && (limit == 0 || capCount < limit) {
		result = append(result, currentSubstring.String())
	}

	if len(result) == 0 {
		return ""
	}
	return strings.Join(result, " ")
}

func handleOffsetAndLimit(data map[string]string) (int, int) {
	offset, err := strconv.Atoi(data["offset"])
	if err != nil {
		loger.Info("invalid offset param, using default value of 0")
		offset = 0
	} else {
		delete(data, "offset")
	}

	limit, err := strconv.Atoi(data["limit"])
	if err != nil {
		loger.Info("invalid limit param, using default value of 10")
		limit = 10
	} else {
		delete(data, "limit")
	}

	return offset, limit
}

func handleTitleAndArtist(data map[string]string) (string, string, int, error) {
	if title, ok := data["title"]; !ok {
		return "", "", http.StatusBadRequest, errors.New("title parameter required")
	} else {
		if artist, ok := data["artist"]; !ok {
			return "", "", http.StatusBadRequest, errors.New("artist parameter required")
		} else {
			return title, artist, http.StatusOK, nil
		}
	}
}
