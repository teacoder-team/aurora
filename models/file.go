package models

type Metadata struct {
	ID     string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Type   string `gorm:"type:varchar(255)" json:"type"`
	Width  int    `gorm:"type:int" json:"width"`
	Height int    `gorm:"type:int" json:"height"`
}

type File struct {
	ID          string   `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Tag         string   `gorm:"type:varchar(255)" json:"tag"`
	Filename    string   `gorm:"type:varchar(255)" json:"filename"`
	Metadata    Metadata `gorm:"foreignKey:MetadataID" json:"metadata"`
	MetadataID  string   `gorm:"type:varchar(255);index" json:"metadata_id"`
	ContentType string   `gorm:"type:varchar(255)" json:"content_type"`
	Size        int      `gorm:"type:int" json:"size"`
	Deleted     *bool    `gorm:"type:boolean" json:"deleted,omitempty"`
	CreatedAt   int64    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   int64    `gorm:"autoUpdateTime" json:"updated_at"`
}
