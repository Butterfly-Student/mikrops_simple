package impl

import (
	"github.com/alijayanet/gembok-backend/internal/domain/entities"
	"gorm.io/gorm"
)

type settingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *settingRepository {
	return &settingRepository{db: db}
}

func (r *settingRepository) Get(key string) (string, error) {
	var s entities.Setting
	if err := r.db.Where("setting_key = ?", key).First(&s).Error; err != nil {
		return "", err
	}
	return s.SettingValue, nil
}

func (r *settingRepository) Set(key, value string) error {
	var s entities.Setting
	result := r.db.Where("setting_key = ?", key).First(&s)
	if result.Error != nil {
		// Create new
		s = entities.Setting{SettingKey: key, SettingValue: value}
		return r.db.Create(&s).Error
	}
	s.SettingValue = value
	return r.db.Save(&s).Error
}

func (r *settingRepository) GetAll() (map[string]string, error) {
	var settings []entities.Setting
	if err := r.db.Find(&settings).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(settings))
	for _, s := range settings {
		result[s.SettingKey] = s.SettingValue
	}
	return result, nil
}
