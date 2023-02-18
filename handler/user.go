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

func (h *Handler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}
	err := h.DB.CreateUser(&user)
	if err != nil {
		h.Logger.Err(err).Msg("could not save user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save user: %s", err.Error())})
	} else {
		c.JSON(http.StatusCreated, gin.H{"user": user})
	}
}

// UpdateUser updates a user's information
// It expects a user ID in the URL and a JSON body that includes the fields to update
// If the user does not exist, it returns a 404
// If the JSON body cannot be parsed, it returns a 400
// If the user cannot be updated in the database, it returns a 500
func (h *Handler) UpdateUser(c *gin.Context) {
	var id int
	var user models.User
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	if err = c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("could not parse request: %s", err.Error())})
		return
	}

	err = h.DB.UpdateUser(id, user)
	if err != nil {
		switch err {
		case db.ErrNoRecord:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find user with id: %d", id)})
		default:
			h.Logger.Err(err).Msg("could not update user")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update user: %s", err.Error())})
		}
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

func (h *Handler) DeleteUser(c *gin.Context) {
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	err = h.DB.DeleteUser(id)
	if err != nil {
		if err == db.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find user with id: %d", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": map[string]string{"message": "user deleted"}})
}

// GetUsers fetches all users from the database and returns them.
// It is used to populate the index page of the app.
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.DB.GetUsers()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch users")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": users})
	}
}

// GetUser returns the user with the given id.
// If no user exists with the given id, it returns a 404 response.
// If there is an error retrieving the user, it returns a 500 response.

func (h *Handler) GetUser(c *gin.Context) {
	var id int
	var err error
	var user models.User
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	user, err = h.DB.GetUserById(id)
	switch err {
	case db.ErrNoRecord:
		log.Printf("could not find user with id: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find user with id: %d", id)})
		return
	case nil:
		c.JSON(http.StatusOK, gin.H{"data": user})
		return
	default:
		log.Printf("error retrieving user: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}
