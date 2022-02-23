package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pet-app/internal/data"
	"pet-app/internal/helper"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	dataFile = "data.json"
)

func (h Handler) showData(c *gin.Context) {
	dataByte, err := helper.ReadFile(dataFile)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	responseApp := data.ResponseApp{}
	err = json.Unmarshal(dataByte, &responseApp)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	b, err := helper.Pretty(responseApp)
	helper.Check(err)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"data": dataByte, "message": "no pretty"})
		return
	}
	c.Writer.Write(b)
	c.Status(http.StatusOK)
}
func (h Handler) addSearchURL(c *gin.Context) {
	inputData := data.SearchURL{}
	err := c.BindJSON(&inputData)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid json")
		return
	}
	if !strings.Contains(inputData.Url, "hh.ru/search/vacancy?") {
		newErrorResponse(c, http.StatusBadRequest, "invalid URL")
		return
	}
	// check response
	client := &http.Client{}
	url := helper.PrepareRequestURL(inputData.Url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cant build request")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "cant sent request")
		return
	}
	if resp.StatusCode != 200 {
		newErrorResponse(c, http.StatusBadRequest, "server on this url not respond")
		return
	}

	inputData.Id = len(h.searchURLs) + 1
	h.searchURLs[inputData.Id] = inputData
	*h.count = 0

	log.Println("[app]New URL added to search options")
	c.JSON(http.StatusOK, gin.H{"message": "add new search url successfully", "url_ID": inputData.Id})
}
func (h Handler) deleteSearchURL(c *gin.Context) {
	inputData := data.DeleteDTOSearchURL{}
	err := c.BindJSON(&inputData)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid json")
		return
	}
	if inputData.Id > len(h.searchURLs) || inputData.Id < 1 {
		newErrorResponse(c, http.StatusBadRequest, "invalid ID")
		return
	}
	delete(h.searchURLs, inputData.Id)
	*h.count = 0
	log.Println("[app]URL with id:", inputData.Id, "has been removed from the search parameters")
	c.JSON(http.StatusOK, gin.H{"message": "search url was deleted"})
}

// auth middlerware
func (h Handler) authorization(tokenAdmin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			newErrorResponse(c, http.StatusUnauthorized, "empty token")
		}
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			newErrorResponse(c, http.StatusUnauthorized, "invalid token")
		}
		if headerParts[1] != tokenAdmin {
			newErrorResponse(c, http.StatusUnauthorized, "wrong token")
		}
	}
}
