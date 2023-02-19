package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PSNAppz/Fold-ELK/db"
	"github.com/PSNAppz/Fold-ELK/models"
	"github.com/gin-gonic/gin"
)

// CreateProject is the handler for creating a new project
func (h *Handler) CreateProject(c *gin.Context) {
	// Bind the JSON request body to a struct
	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}

	// Get the user from the database and verify that the user exists
	user, err := h.DB.GetUserById(req.UserID)
	if err != nil {
		h.Logger.Err(err).Msg("could not get user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not get user: %s", err.Error())})
		return
	}
	// Create the hashtags for the project
	hashtags := make([]models.Hashtag, 0)
	for _, hashtagName := range req.Hashtags {
		hashtag := models.Hashtag{Name: hashtagName}
		tags, err := h.DB.GetOrCreateHashtags([]models.Hashtag{hashtag})
		if err != nil {
			return
		}
		hashtags = append(hashtags, tags[0])
	}
	// Create the project and save it to the database
	project := models.Project{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}
	err = h.DB.CreateProject(&project, &user, &hashtags)
	if err != nil {
		h.Logger.Err(err).Msg("could not save project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not save project: %s", err.Error())})
	} else {
		c.JSON(http.StatusCreated, gin.H{"project": project})
	}
}

func (h *Handler) UpdateProject(c *gin.Context) {
	// Get the project ID from the URL
	var id int
	var req models.CreateProjectRequest
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Err(err).Msg("could not parse request body")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid request body: %s", err.Error())})
		return
	}

	// Get the user from the database and verify that the user exists
	user, err := h.DB.GetUserById(req.UserID)
	if err != nil {
		h.Logger.Err(err).Msg("could not get user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not get user: %s", err.Error())})
		return
	}
	// Create the hashtags for the project
	hashtags := make([]models.Hashtag, 0)
	for _, hashtagName := range req.Hashtags {
		hashtag := models.Hashtag{Name: hashtagName}
		tags, err := h.DB.GetOrCreateHashtags([]models.Hashtag{hashtag})
		if err != nil {
			return
		}
		hashtags = append(hashtags, tags[0])
	}

	project := models.Project{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	// Update the project in the database
	err = h.DB.UpdateProject(id, project, &user, &hashtags)
	if err != nil {
		switch err {
		case db.ErrNoRecord:
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find project with id: %d", id)})
		default:
			h.Logger.Err(err).Msg("could not update project")
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not update project: %s", err.Error())})
		}
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"project": project})
	}
}

// DeleteProject deletes a project.
// The project ID is extracted from the route parameter.
// If the project is successfully deleted, the response will be a
// HTTP 200 with the following body:
//     {
//         "data": {
//             "message": "project deleted"
//         }
//     }
// If the project cannot be found, the response will be a HTTP 404
// with the following body:
//     {
//         "error": "could not find project with id: <project-id>"
//     }
// If there is an error deleting the project, the response will be a
// HTTP 500 with the following body:
//     {
//         "error": <error-message>
//     }

func (h *Handler) DeleteProject(c *gin.Context) {
	var id int
	var err error
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	// Delete user_project associations
	err = h.DB.DeleteUserProjectByProjectId(id)
	if err != nil {
		h.Logger.Err(err).Msg("could not delete user_project")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete user_project: %s", err.Error())})
		return
	}
	// Delete project_hashtag associations
	err = h.DB.DeleteProjectHashtagByProjectId(id)
	if err != nil {
		h.Logger.Err(err).Msg("could not delete project_hashtag")
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("could not delete project_hashtag: %s", err.Error())})
		return
	}

	err = h.DB.DeleteProject(id)
	if err != nil {
		if err == db.ErrNoRecord {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find project with id: %d", id)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": map[string]string{"message": "project deleted"}})
}

// GetProjects fetches all projects from the database and returns them.
// It is used to populate the index page of the app.
func (h *Handler) GetProjects(c *gin.Context) {
	projects, err := h.DB.GetProjects()
	if err != nil {
		h.Logger.Err(err).Msg("Could not fetch projects")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"data": projects})
	}
}

// GetProject returns the project with the given id.
// If no project exists with the given id, it returns a 404 response.
// If there is an error retrieving the project, it returns a 500 response.

func (h *Handler) GetProject(c *gin.Context) {
	var id int
	var err error
	var project models.Project
	if id, err = strconv.Atoi(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}
	project, err = h.DB.GetProjectById(id)
	switch err {
	case db.ErrNoRecord:
		log.Printf("could not find project with id: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("could not find project with id: %d", id)})
		return
	case nil:
		c.JSON(http.StatusOK, gin.H{"data": project})
		return
	default:
		log.Printf("error retrieving project: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func (h *Handler) Search(c *gin.Context) {
	// Get the search query from the request
	var query string
	if query, _ = c.GetQuery("q"); query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no search query present"})
		return
	}
	searchType, err := strconv.Atoi(c.DefaultQuery("type", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be an integer"})
		return
	}
	// Get the fuzziness from the request
	var fuzziness int
	if fuzziness, _ = strconv.Atoi(c.DefaultQuery("fuzziness", "1")); fuzziness < 0 || fuzziness > 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fuzziness must be between 0 and 2"})
		return
	}

	// Create the search query
	// Fuzziness is set to 1 by default, which means that the search will return results
	// that are one edit away from the query string. This can be overridden in the request
	var body string
	switch searchType {
	// Search for a query string in the user_name field
	case 0:
		body = fmt.Sprintf(`{
			"query": {
				"multi_match": {
					"query": "%s",
					"fields": ["user_name"],
					"fuzziness": %d
				}
			}
		}`, query, fuzziness)
	// Search for a query string in the hashtags field
	case 1:
		body = fmt.Sprintf(`{
			"query": {
				"multi_match": {
					"query": "%s",
					"fields": ["hashtags"],
					"fuzziness": %d
				}
			}
		}`, query, fuzziness)
		// Search for a query string in the slug and description fields
	case 2:
		body = fmt.Sprintf(`{
			"query": {
				"multi_match": {
					"query": "%s",
					"fields": ["slug", "description"],
					"fuzziness": %d
				}
			}
		}`, query, fuzziness)

	}

	// Execute the search query
	res, err := h.ESClient.Search(
		h.ESClient.Search.WithContext(context.Background()),
		h.ESClient.Search.WithIndex("projects"),
		h.ESClient.Search.WithBody(strings.NewReader(body)),
		h.ESClient.Search.WithPretty(),
	)
	if err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer res.Body.Close()
	if res.IsError() {
		// Parse the response error into a map
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			h.Logger.Err(err).Msg("error parsing the response body")
		} else {
			// Print the response status and error information.
			h.Logger.Err(fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)).Msg("failed to search fuzzy query")
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": e["error"].(map[string]interface{})["reason"]})
		return
	}
	// Log the response status
	h.Logger.Info().Interface("res", res.Status())

	// Decodes the response body of the Elasticsearch query into the map r
	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		h.Logger.Err(err).Msg("elasticsearch error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": r["hits"]})
}
