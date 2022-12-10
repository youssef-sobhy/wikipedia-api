package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
)

// health godoc
// @Summary		Check if the API is operational.
// @Description	Check if the API is operational.
// @Accept			json
// @Produce		json
// @Success		200	{object}	CheckHealthResponse
// @Failure		500	{object}	InternalServerErrorResponse
// @Router			/api [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, CheckHealthResponse{
		Status: "operational",
	})
}

// search godoc
//
//	@Summary		Search for a short description of a person, place, or thing.
//	@Description	Search for a short description of a person, place, or thing.
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"The name of the person, place, or thing you want to search for."
//	@Success		200		{object}	SuccessResponse
//	@Failure		400		{object}	ErrorResponse
//	@Failure		500		{object}	InternalServerErrorResponse
//	@Router			/api/search [get]
func Search(c *gin.Context) {
	wikipediaURL := os.Getenv("WIKIPEDIA_API_URL")

	if wikipediaURL == "" {
		wikipediaURL = "https://en.wikipedia.org/w/api.php"
	}

	query := c.Query("query")
	if query == "" {
		BadRequestErrorHandler(c, "Query parameter is required.")

		return
	}

	url := fmt.Sprintf("%s?action=query&prop=revisions&titles=%s&rvlimit=1&formatversion=2&format=json&rvprop=content", wikipediaURL, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		InternalServerErrorHandler(c, err)

		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		InternalServerErrorHandler(c, err)

		return
	}

	if resp.StatusCode != http.StatusOK {
		WikipediaApiErrorHandler(c, resp.StatusCode)

		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		InternalServerErrorHandler(c, err)

		return
	}

	var response WikipediaResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		InternalServerErrorHandler(c, err)

		return
	}

	if response.Query.Pages[0].Missing {
		HttpMissingHandler(c)

		return
	}

	content := response.Query.Pages[0].Revisions[0].Content
	re := regexp.MustCompile(`(?mi){{short description\|(.*?)}}`)
	shortDescription := re.FindStringSubmatch(content)

	if len(shortDescription) == 0 {
		HttpNoDescriptionHandler(c)

		return
	}

	HttpSuccessHandler(c, shortDescription[1])
}
