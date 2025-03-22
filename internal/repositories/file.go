package repositories

import (
	"orion/config"
	"orion/internal/models"
	"time"

	"gorm.io/gorm"
)

type FileRepository struct{}

type MetadataRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func NewMetadataRepository() *MetadataRepository {
	return &MetadataRepository{}
}

func (r *FileRepository) GetFileByID(id, tag string) (*models.File, error) {
	var file models.File
	if err := config.DB.First(&file, "id = ? AND tag = ? AND deleted = ?", id, tag, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &file, nil
}

func (r *FileRepository) CreateFile(file *models.File) error {
	if err := config.DB.Create(file).Error; err != nil {
		return err
	}
	return nil
}

func (r *FileRepository) SoftDeleteFile(id, tag string) error {
	var file models.File

	if err := config.DB.First(&file, "id = ? AND tag = ? AND deleted = ?", id, tag, false).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	currentTime := time.Now()
	file.Deleted = new(bool)
	*file.Deleted = true
	file.DeletedAt = &currentTime

	if err := config.DB.Save(&file).Error; err != nil {
		return err
	}

	return nil
}
