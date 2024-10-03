package main

import (
	httpH "choomandco/efimobi/http_handler"
	models "choomandco/efimobi/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"gorm.io/gorm"
)

func getSoundData(db *gorm.DB, data map[string]string) (int, string, error) {
	if title, ok := data["title"]; !ok {
		return http.StatusBadRequest, "", errors.New("title field required")
	} else {
		artist := data["artist"]
		offset, limit := handleOffsetAndLimit(data)

		soundData, err := models.SelectSoundData(db, title, artist, offset, limit)
		if err != nil {
			return http.StatusInternalServerError, "", err
		} else {
			return http.StatusOK, soundData, nil
		}
	}
}

func getArtistData(db *gorm.DB, data map[string]string) (int, string, error) {
	offset, limit := handleOffsetAndLimit(data)

	artistData, err := models.SelectArtistData(db, data["artist"], offset, limit)
	if err != nil {
		return http.StatusInternalServerError, "", err
	} else {
		res, err := httpH.ToJSON(artistData)
		if err != nil {
			return http.StatusInternalServerError, "", err
		} else {
			return http.StatusOK, res, nil
		}
	}
}

func getAllSoundsData(db *gorm.DB, data map[string]string) (int, string, error) {
	offset, limit := handleOffsetAndLimit(data)
	wholeLib, err := models.SelectWholeLibData(db, data, offset, limit)
	if err != nil {
		return http.StatusInternalServerError, "", err
	} else {
		return http.StatusOK, wholeLib, nil
	}
}

func patchDeleteSound(db *gorm.DB, data map[string]string) (int, string, error) {
	title, artist, code, err := handleTitleAndArtist(data)
	if err != nil {
		loger.Warn(err.Error())
		return code, "", err
	} else {
		code, resp, err := models.DeleteSoundFromLib(db, title, artist)
		return code, resp, err
	}
}

func postNewSound(db *gorm.DB, data map[string]string) (int, string, error) {
	title, artist, code, err := handleTitleAndArtist(data)
	if err != nil {
		loger.Warn(err.Error())
		return code, "", err
	} else {
		apiURL := fmt.Sprintf("%s/search_track?title=%s&artist=%s", tokenAddress, url.QueryEscape(title), url.QueryEscape(artist))

		if _, resp, err := httpH.CallOuterApi(apiURL); err != nil {
			loger.Warn(err)
			resp, err := models.CreateSoundSimple(db, title, artist)
			if err != nil {
				return http.StatusInternalServerError, "", err
			} else {
				return http.StatusOK, resp, nil
			}
		} else {
			res, err := models.CreateSoundWithJSON(db, resp)
			if err != nil {
				resp, err := models.CreateSoundSimple(db, title, artist)
				if err != nil {
					return http.StatusInternalServerError, "", err
				} else {
					return http.StatusCreated, resp, nil
				}
			} else {
				return http.StatusCreated, res, nil
			}
		}
	}
}

func getSoundLyrics(db *gorm.DB, data map[string]string) (int, string, error) {
	title, artist, code, err := handleTitleAndArtist(data)
	if err != nil {
		return code, "", err
	} else {
		limit, offset := handleOffsetAndLimit(data)
		if lyr, err := models.SelectLyrics(db, title, artist); err == nil {
			res := extractStrings(lyr, offset, limit)
			return http.StatusOK, res, nil
		} else {
			return http.StatusInternalServerError, "", err
		}
	}
}

func patchExistingSound(db *gorm.DB, data map[string]string) (int, string, error) {
	title, artist, code, err := handleTitleAndArtist(data)
	if err != nil {
		loger.Warn(err)
		return code, "", err
	} else {
		if code, resp, err := models.UpdateExistingSound(db, title, artist, data); err != nil {
			return code, "", err
		} else {
			return code, resp, nil
		}
	}
}

func executeSQLFile(db *gorm.DB, filepath string) string {
	file, err := os.Open(filepath)
	if err != nil {
		loger.Warn("failed to open SQL file")
		return ""
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		loger.Warn("failed to read SQL file")
		return ""
	}

	err = db.Exec(string(content)).Error
	if err != nil {
		loger.Warn("failed to execute SQL content")
		return ""
	}

	nullTable = false
	resp := "dml executed successfully"
	return resp
}
