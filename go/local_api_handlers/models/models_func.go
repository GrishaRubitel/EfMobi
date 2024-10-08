package models

import (
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var loger = logrus.New()

// Миграция моделей в БД
func MigrateModels(db *gorm.DB) {
	if err := db.Migrator().CreateTable(ArtistData{}); err != nil {
		loger.Warn("failed to create artist_data table: " + err.Error())
	}

	if err := db.Migrator().CreateTable(SoundData{}); err != nil {
		loger.Warn("failed to create sound_data table: " + err.Error())
	}

	createForeignKey(db)
}

// Замена пробелов на "любое количество символов" для фильтрации во время выполнения SQL запросов
func skipSpaces(str string) string {
	return "%" + strings.ReplaceAll(str, " ", "%") + "%"
}

// Создание внешнего ключа между двумя существующими таблицами
func createForeignKey(db *gorm.DB) {
	err := db.Exec(`ALTER TABLE sound_data
                    ADD CONSTRAINT fk_sound_data_artist_data
                    FOREIGN KEY (artist) REFERENCES artist_data(artist) ON DELETE CASCADE ON UPDATE CASCADE;`).Error
	if err != nil {
		loger.Warn("failed to create foreign key: " + err.Error())
	}
}
