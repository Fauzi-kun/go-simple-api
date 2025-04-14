package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var validate = validator.New()

func main() {
	initDB()
	r := gin.Default()

	r.GET("/notes",get)
	r.POST("/notes",post)
	r.PUT("/notes/:id",put)
	r.DELETE("/notes/:id",delete)
	r.POST("/login", login)
	r.GET("/profile",AuthMiddleware(),profile)
	r.Run(":8080")
}
func get(c *gin.Context){
	var allNotes []Note
	db.Find(&allNotes)
	c.JSON(http.StatusOK, allNotes)
}
func post(c *gin.Context){
	var newNote Note
	if err := c.BindJSON(&newNote); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Format JSON salah!"})
		return
	}

	//Validate
	if err := validate.Struct(newNote); err != nil{
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

	//bind data baru
	var updatedData Note
	if err := c.BindJSON(&updatedData); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": "Format JSON salah!"})
		return
	}

	//Validate
	if err := validate.Struct(updatedData); err != nil{
		c.JSON(http.StatusBadRequest,gin.H{"error": err.Error()})
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
func login(c *gin.Context){
	var body struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&body);err != nil{
		c.JSON(400, gin.H{"error": "Format salah"})
		return
	}
	if body.Username != "admin" || body.Password != "rahasia"{
		c.JSON(401,gin.H{"error": "Username/password salah"})
		return
	}
	token, _ := GenerateToken(body.Username)
	c.JSON(200, gin.H{"Token": token})
}
func profile(c *gin.Context){
	username := c.MustGet("username").(string)
	c.JSON(200, gin.H{"message": "Hai " + username})
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
