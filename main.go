package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var notes = []Note{
	{ID: 1, Title: "Belajar Go", Content: "Pelajari Gin dan GORM"},
}
var db *gorm.DB

func main() {
	initDB()
	r := gin.Default()

	r.GET("/notes",get)
	r.POST("/notes",post)
	r.PUT("/notes/:id",put)
	r.DELETE("/notes/:id",delete)
	r.Run(":8080")
}
func get(c *gin.Context){
	var allNotes []Note
	db.Find(&allNotes)
	c.JSON(http.StatusOK, allNotes)
}
func post(c *gin.Context){
	var newNote Note
	if err := c.ShouldBindJSON(&newNote); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
		return
	} 
	db.Create(&newNote)
	c.JSON(http.StatusCreated,newNote)
}
func put(c *gin.Context){
	id := c.Param("id")
	var note Note

	//cari note
	if err := db.First(&note,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error": "Note nggak ditemukan!"})
		return
	}

	//bin data baru
	var updatedData Note
	if err := c.BindJSON(&updatedData); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Format JSON salah!"})
		return
	}
	
	//update field
	note.Title = updatedData.Title
	note.Content = updatedData.Content
	db.Save(&note)
	c.JSON(http.StatusOK,note)
}
func delete(c *gin.Context){
	id := c.Param("id")
	var note Note

	//Cari note
	if err := db.First(&note,id).Error; err != nil{
		c.JSON(http.StatusNotFound,gin.H{"error": "Note tidak ditemukan!"})
		return
	}
	db.Delete(&note)
	c.JSON(http.StatusOK,gin.H{"Message": "Note berhasil dihapus!"})
}
func initDB(){
	dsn := "host=localhost user=postgres password=fauzi dbname=notes_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil{
		panic("Gagal koneksi ke database")
	}
	db = database
	db.AutoMigrate(&Note{})
}