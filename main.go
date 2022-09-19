package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	Id       string `gorm: "primary_key" json: "id" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
	City     string `json:"city" validate:"required"`
}

func getAllUsers(c *gin.Context) {
	var allusers = []Users{}
	err := DB.Find(&allusers).Error
	if err != nil {
		c.IndentedJSON(http.StatusNoContent, err.Error())
		return
	}
	c.IndentedJSON(http.StatusFound, allusers)

}

func getUserByID(c *gin.Context) {
	var user Users
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "id cannot be empty"})
		return
	}

	if err := DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "could not find the requested user"})
		return
	}
	c.IndentedJSON(http.StatusFound, gin.H{"User found :)": user})

}

func createNewUser(c *gin.Context) {
	var newuser Users

	if err := c.ShouldBindJSON(&newuser); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"error-BindJSON": err.Error()})
		return
	}

	validator := validator.New()
	if err := validator.Struct(&newuser); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"Validation-error": err.Error()})
		return
	}

	err := DB.Create(&newuser).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Error creating User :( ": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"New user created :)": newuser})

}

func deleteUser(c *gin.Context) {
	var userToDelete Users

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "id cannot be empty"})
		return
	}

	err := DB.Delete(userToDelete, id).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "could not delete user", "error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, gin.H{"message": "Delete successful."})

}

func UpdateById(c *gin.Context) {
	var updateuser Users

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Id should not be empty"})
		return
	}

	if err := c.ShouldBindJSON(&updateuser); err != nil {
		c.IndentedJSON(http.StatusUnprocessableEntity, gin.H{"not valid": err.Error()})
		return
	}
	err := DB.Where("id = ?", id).Updates(updateuser).Error
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"if found": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusAccepted, gin.H{"updated succefully": updateuser})

}

var DB *gorm.DB

func setupDB() {
	connectionString := "host=localhost user=postgres password=12345 dbname=TestUsers port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatal("could not connect to database")
		return
	}
	DB = db

}

func migrate() {
	if (DB.Migrator().CurrentDatabase() != "") && (DB.Migrator().HasTable(&Users{})) {
		return
	}
	err := DB.AutoMigrate(&Users{})
	if err != nil {
		log.Fatal("could not migrate to db")
		return
	}
}

func main() {
	setupDB()
	migrate()
	router := gin.Default()
	router.GET("/allusers", getAllUsers)
	router.GET("/user/:id", getUserByID)
	router.POST("/newuser", createNewUser)
	router.DELETE("/delete/:id", deleteUser)
	router.PUT("/update/:id", UpdateById)
	router.Run("localhost:8080")

}
