package image

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
)

const (
	UserAvatarImage      Type = "user_avatar"
	EventIconImage       Type = "event_icon"
	AchievementIconImage Type = "achievement_icon"
)

type Type string

type Options struct {
	Users       ValidationOptions `yaml:"users"`
	Events      ValidationOptions `yaml:"events"`
	Achievement ValidationOptions `yaml:"achievements"`
}

type ValidationOptions struct {
	MaxSize      int64    `yaml:"max_size"`
	ContentTypes []string `yaml:"content_types"`
}

type Service interface {
	Validate(ctx context.Context, imageType Type, image *graphql.Upload) error
}

type service struct {
	validationMap map[Type]ValidationOptions
}

func (s *service) Validate(ctx context.Context, imageType Type, image *graphql.Upload) error {
	validationRules, ok := s.validationMap[imageType]
	if !ok {
		return fmt.Errorf("invalid image type")
	}
	if image.Size > validationRules.MaxSize {
		return fmt.Errorf("image to large")
	}
	if !contains(validationRules.ContentTypes, image.ContentType) {
		return fmt.Errorf("wrong content type")
	}
	return nil
}

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func GenerateFilename(image *graphql.Upload) string {
	ext := filepath.Ext(image.Filename)
	generatedName := uuid.New().String()
	return generatedName + ext
}

func New(opts Options) Service {
	validationMap := make(map[Type]ValidationOptions, 3)
	validationMap[UserAvatarImage] = opts.Users
	validationMap[EventIconImage] = opts.Events
	validationMap[AchievementIconImage] = opts.Achievement
	return &service{
		validationMap: validationMap,
	}
}
