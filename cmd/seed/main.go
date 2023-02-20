package main

import (
	"math/rand"
	"os"
	"strings"

	"github.com/PSNAppz/Fold-ELK/db"
	"github.com/PSNAppz/Fold-ELK/models"
	"github.com/bxcodec/faker/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	dbConfig := db.Config{
		// hard-coding db connection info since we only have to seed on the dev machine
		Host:     "localhost",
		Port:     5432,
		Username: "fold-elk",
		Password: "password",
		DbName:   "fold_elk",
	}

	dbInstance, err := db.Init(dbConfig)
	if err != nil {
		logger.Err(err).Msg("Connection failed")
		os.Exit(1)
	}
	logger.Info().Msg("Connected to Database")

	// Seed dummy users
	logger.Info().Msg("Seeding Users to Database")
	for i := 0; i < 15; i++ {
		user := &models.User{
			Name: faker.Name(),
		}
		err = dbInstance.CreateUser(user)
		if err != nil {
			log.Err(err).Msg("failed to save record")
		}
	}
	logger.Info().Msg("Completed Seeding Users to Database")

	// Create Hashtags
	logger.Info().Msg("Seeding Hashtags to Database")
	for i := 0; i < 15; i++ {
		ht := &models.Hashtag{
			Name: faker.Word(),
		}
		err = dbInstance.CreateHashtag(ht)
		if err != nil {
			log.Err(err).Msg("failed to save record")
		}
	}
	logger.Info().Msg("Completed Seeding Hashtags to Database")

	// Seed dummy projects
	logger.Info().Msg("Seeding Projects to Database")
	for i := 0; i < 10; i++ {
		name := faker.Sentence()
		project := &models.Project{
			Name:        name,
			Slug:        strings.ReplaceAll(strings.ToLower(name), " ", "-"),
			Description: faker.Paragraph(),
		}

		// get random user
		user, err := dbInstance.GetUserById(rand.Intn(15) + 1)
		if err != nil {
			log.Err(err).Msg("failed to get user")
		}

		// get random hashtags
		hashtagIds := make([]models.Hashtag, 0)
		for i := 0; i < 3; i++ {
			hashtagIds, err = dbInstance.GetHashtags()
			if err != nil {
				log.Err(err).Msg("failed to get hashtags")
			}
		}
		// create project
		err = dbInstance.CreateProject(project, &user, &hashtagIds)
		if err != nil {
			log.Err(err).Msg("failed to save record")
		}
	}
	logger.Info().Msg("Completed Seeding Projects to Database")

}
