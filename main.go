package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var notes = []Note{
	{ID: 1, Title: "Belajar Go", Content: "Pelajari Gin dan GORM"},
}

func main() {
	r := gin.Default()

	r.GET("/notes",func(c *gin.Context){
		c.JSON(http.StatusOK, notes)
	})
	r.POST("/notes",post)
	r.PUT("/notes/:id",put)
	r.Run(":8080")
}
func post(c *gin.Context){
	var newNote Note
	if err := c.ShouldBindJSON(&newNote); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	} 
	newNote.ID = len(notes) + 1
	notes = append(notes, newNote)
	c.JSON(http.StatusCreated,newNote)
}
func put(c *gin.Context){
	idParam := c.Param("id")
	id, err :=strconv.Atoi(idParam)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "ID nggak valid!"})
		return
	}
	var updatedNote Note
	if err := c.BindJSON(&updatedNote); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Format JSON salah!"})
		return
	}
	for i, note := range notes{
		if note.ID == id{
			notes[i].Title = updatedNote.Title
			notes[i].Content = updatedNote.Content
			c.JSON(http.StatusOK,notes[i])
			return
		}
	}
	c.JSON(http.StatusNotFound,gin.H{"error": "Note nggak ditemukan!"})
}