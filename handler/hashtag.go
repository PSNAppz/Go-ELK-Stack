package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/PSNAppz/Fold-ELK/db"
	"github.com/PSNAppz/Fold-ELK/models"
	"github.com/gin-gonic/gin"
)

// CreateHashtag creates a new hashtag in the database.
func (h *Handler) CreateHashtag(c *gin.Context) {
	// Bind the request body to a new Hashtag struct
	var hashtag models.Hashtag
	if err := c.ShouldBindJSON(&hashtag); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}
	// Attempt to create the new hashtag in the database
	err := h.DB.CreateHashtag(&hashtag)
	if err != nil {
		h.Logger.Err(err).Msg("could not save hashtag")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save hashtag: %s", err.Error())})
	} else {
		// If successful, return the new hashtag in the response with 201 Created status code
		c.JSON(http.StatusCreated, gin.H{"hashtag": hashtag})
	}
}

// UpdateHashtag updates a hashtag's information
func (h *Handler) UpdateHashtag(c *gin.Context) {
	var id int
	var hashtag models.Hashtag
	var err error
	// It expects a hashtag ID in the URL and a JSON body that includes the fields to update
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hashtag id"})
		return
	}
	if err = c.ShouldBindJSON(&hashtag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse request: %s", err.Error())})
		return
	}
	// Update the hashtag in the database
	err = h.DB.UpdateHashtag(id, hashtag)
	if err != nil {
		switch err {
		case db.ErrNoRecord:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find hashtag with id: %d", id)})
		default:
			h.Logger.Err(err).Msg("could not update hashtag")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update hashtag: %s", err.Error())})
		}
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"hashtag": hashtag})
	}
}

// Delete a hashtag from the database
func (h *Handler) DeleteHashtag(c *gin.Context) {
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hashtag id"})
		return
	}
	// Delete all hashtag_project associations
	err = h.DB.DeleteHashtagProjectByHashtagId(id)
	if err != nil {
		h.Logger.Err(err).Msg("could not delete hashtag_project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete hashtag_project: %s", err.Error())})
		return
	}
	err = h.DB.DeleteHashtag(id)
	if err != nil {
		if err == db.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find hashtag with id: %d", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": map[string]string{"message": "hashtag deleted"}})
}

// GetHashtags fetches all hashtags from the database and returns them.
// It is used to populate the index page of the app.
func (h *Handler) GetHashtags(c *gin.Context) {
	hashtags, err := h.DB.GetHashtags()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch hashtags")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": hashtags})
	}
}

// GetHashtag returns the hashtag with the given id.
// If no hashtag exists with the given id, it returns a 404 response.
// If there is an error retrieving the hashtag, it returns a 500 response.

func (h *Handler) GetHashtag(c *gin.Context) {
	var id int
	var err error
	var hashtag models.Hashtag
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hashtag id"})
		return
	}
	hashtag, err = h.DB.GetHashtagById(id)
	switch err {
	case db.ErrNoRecord:
		log.Printf("could not find hashtag with id: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find hashtag with id: %d", id)})
		return
	case nil:
		c.JSON(http.StatusOK, gin.H{"data": hashtag})
		return
	default:
		log.Printf("error retrieving hashtag: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
