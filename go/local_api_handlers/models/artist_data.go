package models

import (
	httpH "choomandco/efimobi/http_handler"
	"errors"
	"log"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ArtistData struct {
	Artist      string         `gorm:"primaryKey;size:255" json:"artist"`
	SpotifyID   string         `gorm:"size:255" json:"spotify_id"`
	SpotifyLink string         `gorm:"size:255" json:"spotify_link"`
	YoutubeLink string         `gorm:"size:255" json:"youtube_link"`
	Genres      datatypes.JSON `json:"genres"`
}

func SelectArtistData(db *gorm.DB, artist string, offset int, limit int) ([]ArtistData, error) {
	var artistData []ArtistData

	query := db.Table("artist_data ad").Limit(limit).Offset(offset)
	if artist != "" {
		query.Where("artist ilike ?", skipSpaces(artist))
	}

	result := query.Find(&artistData)
	if err := result.Error; err != nil {
		log.Panic("Error while selecting artists - " + err.Error())
		return []ArtistData{}, errors.New("error while selecting artists")
	} else {
		return artistData, nil
	}
}

func CreateArtistWithJSON(db *gorm.DB, in map[string]interface{}) (string, error) {
	var newArtist ArtistData

	newArtist.Artist = in["name"].(string)
	newArtist.SpotifyID = in["id"].(string)
	newArtist.SpotifyLink = in["external_urls"].(map[string]interface{})["spotify"].(string)

	result := db.Table("artist_data").Create(&newArtist)
	if result.Error != nil {
		loger.Warn("error while adding new artist - ", result.Error)
		return "", errors.New("errors while adding new artist")
	} else {
		if res, err := httpH.ToJSON(newArtist); err != nil {
			return "", err
		} else {
			return res, nil
		}
	}
}

func CreateArtistSimple(db *gorm.DB, artist string) {
	var newArtist ArtistData
	newArtist.Artist = artist

	result := db.Table("artist_data").Create(&newArtist)
	if result.Error != nil {
		loger.Warn("error while adding new artist - ", result.Error)
		return
	}
}
