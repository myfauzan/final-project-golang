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

func CreateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	SocialMedia := models.SocialMedia{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID

	err := db.Debug().Create(&SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": SocialMedia.ID,
		"name": SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id": SocialMedia.UserID,
		"created_at": SocialMedia.CreatedAt,
	})
}


func GetSocialMedia(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    SocialMedia := []models.SocialMedia{}

    if err := db.Preload("User").Where("user_id = ?", userID).Find(&SocialMedia).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to retrieve SocialMedia",
        })
        return
    }


	socialMedias := []gin.H{}


	for _, socialmedia := range SocialMedia {
    	socialMedias = append(socialMedias, gin.H{
        	"id": socialmedia.ID,
			"name": socialmedia.Name,
			"social_media_url": socialmedia.SocialMediaUrl,
			"user_id": socialmedia.UserID,
			"created_at": socialmedia.CreatedAt,
			"User": gin.H{
				"id": socialmedia.User.ID,
				"username": socialmedia.User.UserName,
				"email": socialmedia.User.Email,
			},
    	})
	}


	response := gin.H{"social_medias": socialMedias}


	c.JSON(http.StatusOK, response)

}


func UpdateSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.UserID = userID
	SocialMedia.ID = uint(socialMediaId)

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaId).Updates(models.SocialMedia{Name: SocialMedia.Name, SocialMediaUrl: SocialMedia.SocialMediaUrl}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": SocialMedia.ID,
		"name": SocialMedia.Name,
		"social_media_url": SocialMedia.SocialMediaUrl,
		"user_id": SocialMedia.UserID,
		"updated_at": SocialMedia.UpdatedAt,
	})
}

func DeleteSocialMedia(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	SocialMedia := models.SocialMedia{}

	socialMediaId, _ := strconv.Atoi(c.Param("socialMediaId"))
	id := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&SocialMedia)
	} else {
		c.ShouldBind(&SocialMedia)
	}

	SocialMedia.ID = uint(id)

	err := db.Model(&SocialMedia).Where("id = ?", socialMediaId).Delete(SocialMedia).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"messege": "Your social media has been succesfully deleted",
	})
}