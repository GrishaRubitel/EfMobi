package models

import (
	httpH "choomandco/efimobi/http_handler"

	"errors"
	"fmt"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type SoundData struct {
	ID          uint      `gorm:"primaryKey"`
	Artist      string    `gorm:"size:255;index" json:"artist"`
	Title       string    `gorm:"size:255;index" json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Lyrics      string    `gorm:"type:text" json:"lyrics"`
	SpotifyLink string    `gorm:"unique;size:255" json:"spotify_link"`
	SpotifyID   string    `gorm:"unique;size:255" json:"spotify_id"`
	VideoLink   string    `gorm:"unique;size:255" json:"video_link"`
}

// Удаление записи о треке по названию и имени исполнителя
func DeleteSoundFromLib(db *gorm.DB, title string, artist string) (int, string, error) {
	var sound SoundData
	if err := db.Where("title ilike ? AND artist ilike ?", skipSpaces(title), skipSpaces(artist)).First(&sound).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			loger.Warn("no requested sound records found; title - ", title, "; artist - ", artist)
			return http.StatusNotFound, "", errors.New("no records found")
		}
		loger.Warn("error while searchin for sound record; title - ", title, "; artist - ", artist)
		return http.StatusBadRequest, "", errors.New("error while wearching for record")
	}

	if err := db.Delete(&sound).Error; err != nil {
		loger.Warn("error while deleting sound record - ", err.Error())
		return http.StatusInternalServerError, "", errors.New("error while deleting sound record")
	}

	resp, _ := httpH.ToJSON(sound)
	loger.Info("sound record deleted - ", resp)
	return http.StatusOK, "record successfully deleted - " + resp, nil
}

// Извлечение информации об одном или нексольких треках с одним и тем же названием
func SelectSoundData(db *gorm.DB, title string, artist string, offset int, limit int) (string, error) {
	var soundData []SoundData

	query := db.Table("sound_data sd").Where("title ilike ?", skipSpaces(title)).Offset(offset).Limit(limit)
	if artist != "" {
		query.Where("artist ilike ?", artist)
	}

	result := query.Find(&soundData)
	if err := result.Error; err != nil {
		loger.Panic("error while selecting sounds - ", err.Error())
		return "", errors.New("error while selecting sounds")
	} else {
		json, err := httpH.ToJSON(soundData)
		if err != nil {
			return "", err
		}
		return json, nil
	}
}

// Извлечение информации обо всех треках библиотеки с фильтрацией по всем (почти) атрибутам
func SelectWholeLibData(db *gorm.DB, data map[string]string, offset int, limit int) (string, error) {
	var libData []SoundData

	query := db.Table("sound_data sd").Offset(offset).Limit(limit)

	// if da, ok := data["release_date"]; ok {
	// 	d, err := time.Parse("02-01-2006", da)
	// 	if err != nil {
	// 		loger.Warn("invalid or non-existing date variable, skipping it - ", err.Error())
	// 	} else {
	// 		query = query.Where("release_date = ?", d.Format("2006-01-02"))
	// 		delete(data, "release_date")
	// 	}
	// }

	for key, val := range data {
		filter := fmt.Sprintf("%s ilike ?", key)
		query = query.Where(filter, skipSpaces(val))
	}

	result := query.Find(&libData)
	if err := result.Error; err != nil {
		loger.Panic("error while selecting lib - ", err.Error())
		return "", errors.New("error while selecting library")
	} else {
		json, err := httpH.ToJSON(libData)
		if err != nil {
			return "", err
		}
		return json, nil
	}
}

// Создание записи о треке с использованием информации, полученной из Spotify Web API
func CreateSoundWithJSON(db *gorm.DB, in map[string]interface{}) (string, error) {
	var newSound SoundData

	tracks, ok := in["tracks"].(map[string]interface{})
	if !ok {
		loger.Warn("error while selecting 'tracks' map")
		return "", errors.New("error while selecting 'tracks' map")
	}

	items, ok := tracks["items"].([]interface{})
	if !ok {
		loger.Warn("error while selecting 'items' map")
		return "", errors.New("error while selecting 'items' map")
	}

	if len(items) > 0 {
		firstItem := items[0].(map[string]interface{})

		artists, ok := firstItem["artists"].([]interface{})
		if !ok {
			loger.Warn("error while selecting 'artists' map")
			return "", errors.New("error while selecting 'artists' map")
		}

		if len(artists) > 0 {
			firstArtist := artists[0].(map[string]interface{})

			_, _ = CreateArtistWithJSON(db, firstArtist)

			newSound.Artist = firstArtist["name"].(string)
		}

		newSound.Title = firstItem["name"].(string)
		//newSound.ReleaseDate = firstItem["release_date"].(time.Time)
		newSound.SpotifyID = firstItem["id"].(string)
		newSound.SpotifyLink = firstItem["external_urls"].(map[string]interface{})["spotify"].(string)

	} else {
		loger.Warn("map items is empty")
		return "", errors.New("map items is empty")
	}

	result := db.Table("sound_data").Create(&newSound)
	if result.Error != nil {
		er := "error while adding new sound record"
		loger.Panic(er)
		return "", errors.New(er)
	} else {
		if ret, err := httpH.ToJSON(newSound); err != nil {
			return "", err
		} else {
			return ret, nil
		}
	}
}

// Извлечение текста трека с пагинацией по строкам (куплетам)
func SelectLyrics(db *gorm.DB, title string, artist string) (string, error) {
	var sound SoundData
	res := db.Table("sound_data").Where("title ilike ?", skipSpaces(title)).Where("artist ilike ?", skipSpaces(artist)).Find(&sound)
	if res.Error != nil {
		loger.Warn("error while selecting lyrics - ", res.Error.Error())
		return "", res.Error
	} else {
		return sound.Lyrics, nil
	}
}

// Создание записи о треке с использованием информации, полученной от клиента
func CreateSoundSimple(db *gorm.DB, title string, artist string) (string, error) {
	var sound SoundData
	sound.Artist = artist
	sound.Title = title

	CreateArtistSimple(db, artist)

	res, err := httpH.ToJSON(sound)
	return res, err
}

// Обновление информации о существующем треке
func UpdateExistingSound(db *gorm.DB, title string, artist string, data map[string]string) (int, string, error) {
	var soundData SoundData

	query := db.Table("sound_data sd").Where("title ilike ?", skipSpaces(title)).Where("artist ilike ?", skipSpaces(artist)).First(&soundData)
	if query.Error != nil {
		loger.Warn("error while selecting sounds - ", query.Error.Error())
		return http.StatusNotFound, "", errors.New("error while selecting sounds")
	}

	if date, ok := data["release_date"]; ok {
		d, err := time.Parse("02-01-2006", date)
		if err != nil {
			loger.Warn("invalid or non-existing date variable, skipping it - ", err.Error())
		} else {
			soundData.ReleaseDate = d
		}
	}

	if lyrics, ok := data["lyrics"]; ok {
		soundData.Lyrics = lyrics
	}

	if spLink, ok := data["spotify_link"]; ok {
		soundData.SpotifyLink = spLink
	}

	if spId, ok := data["spotify_id"]; ok {
		soundData.SpotifyID = spId
	}

	if vidLink, ok := data["video_link"]; ok {
		soundData.VideoLink = vidLink
	}

	if err := db.Table("sound_data").Save(&soundData).Error; err != nil {
		loger.Warn("error while updating record - ", err)
		return http.StatusInternalServerError, "", errors.New("error while updating record")
	} else {
		if res, err := httpH.ToJSON(soundData); err != nil {
			return http.StatusInternalServerError, "", err
		} else {
			return http.StatusAccepted, res, nil
		}
	}
}
