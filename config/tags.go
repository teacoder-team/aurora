package config

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type ContentType string

const (
	ContentTypeImage ContentType = "Image"
	ContentTypeVideo ContentType = "Video"
	ContentTypeAudio ContentType = "Audio"
)

type Tag struct {
	MaxSize             int          `toml:"max_size"`
	UseUlid             bool         `toml:"use_ulid,omitempty"`
	Enabled             bool         `toml:"enabled,omitempty"`
	ServeIfFieldPresent []string     `toml:"serve_if_field_present,omitempty"`
	RestrictContentType *ContentType `toml:"restrict_content_type,omitempty"`
}

type ConfigTags struct {
	Tags        map[string]Tag `toml:"tags"`
	JpegQuality int            `toml:"jpeg_quality"`
}

var configTags *ConfigTags

func LoadConfigTags() error {
	file, err := os.Open("config.toml")
	if err != nil {
		log.Fatalf("❌ Unable to open tags config file: %v\n", err)
		return err
	}
	defer file.Close()

	if configTags == nil {
		configTags = &ConfigTags{}
	}

	log.Println("🔄 Loading config from config.toml...")

	if _, err := toml.DecodeReader(file, configTags); err != nil {
		log.Printf("❌ Unable to parse tags config file: %v\n", err)
		return err
	}

	log.Println("✅ Config loaded successfully!")

	return nil
}

func GetConfigTags() (*ConfigTags, error) {
	if configTags == nil {
		log.Fatalf("❌ Config tags not loaded!")
		return nil, fmt.Errorf("config tags not loaded")
	}

	log.Println("🔑 Returning loaded config tags.")
	return configTags, nil
}
