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

// Функция API для получения информации об одном или нескольких треков с одним и тем же названием
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

// Функция API для получения информации об указанном артисте, либо же обо всех с базе
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

// Функция API для получения целой библотеки треков, с фильтром по всем атрибутам
func getAllSoundsData(db *gorm.DB, data map[string]string) (int, string, error) {
	offset, limit := handleOffsetAndLimit(data)
	wholeLib, err := models.SelectWholeLibData(db, data, offset, limit)
	if err != nil {
		return http.StatusInternalServerError, "", err
	} else {
		return http.StatusOK, wholeLib, nil
	}
}

// Функция API для удаления трека с соответсвующим названием соответствующего исполнителя
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

// Функция API для добавления в базу нового трека
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

// Функция API для получения текста определенного трека, будь то всего, будь то с ограничением по количеству строк (куплетов)
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

// Функция API для обновления информации о существующем треке
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

// Функция API для "запуска" DML скрипта (заполнение БД тестовыми данными)
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

	nullTable = false

	err = db.Exec(string(content)).Error
	if err != nil {
		loger.Warn("failed to execute SQL content")
		return ""
	}

	resp := "dml executed successfully"
	return resp
}
