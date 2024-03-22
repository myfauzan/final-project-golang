package controllers

import (
	"final-project/database"
	"final-project/helpers"
	"final-project/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Photo := models.Photo{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = userID

	err := db.Debug().Create(&Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": Photo.ID,
		"title": Photo.Title,
		"caption": Photo.Caption,
		"photo_url": Photo.PhotoUrl,
		"user_id": Photo.UserID,
		"created_at": Photo.CreatedAt,
	})
}


func GetPhoto(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    Photo := []models.Photo{}

    if err := db.Preload("User").Where("user_id = ?", userID).Find(&Photo).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to retrieve Photo",
        })
        return
    }

    var response []gin.H
    for _, photo := range Photo {
        response = append(response, gin.H{
            "id":         photo.ID,
            "title":      photo.Title,
            "caption":    photo.Caption,
            "photo_url":  photo.PhotoUrl,
            "user_id":    photo.UserID,
            "created_at": photo.CreatedAt,
            "updated_at": photo.UpdatedAt,
            "User": gin.H{
                "email":    photo.User.Email,
                "username": photo.User.UserName,
            },
        })
    }

    c.JSON(http.StatusOK, response)
}


func UpdatePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.UserID = userID
	Photo.ID = uint(photoId)

	err := db.Model(&Photo).Where("id = ?", photoId).Updates(models.Photo{Title: Photo.Title, Caption: Photo.Caption, PhotoUrl: Photo.PhotoUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": Photo.ID,
		"title": Photo.Title,
		"caption": Photo.Caption,
		"photo_url": Photo.PhotoUrl,
		"user_id": Photo.UserID,
		"updated_at": Photo.UpdatedAt,
	})
}

func DeletePhoto(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Photo := models.Photo{}

	photoId, _ := strconv.Atoi(c.Param("photoId"))
	id := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Photo)
	} else {
		c.ShouldBind(&Photo)
	}

	Photo.ID = uint(id)

	err := db.Model(&Photo).Where("id = ?", photoId).Delete(Photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "Your photo has been succesfully deleted",
	})
}