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

func CreateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)

	Comment := models.Comment{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID

	err := db.Debug().Create(&Comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": Comment.ID,
		"message": Comment.Message,
		"photo_id": Comment.PhotoID,
		"user_id": Comment.UserID,
		"created_at": Comment.CreatedAt,
	})
}

func GetComent(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))

    Comment := []models.Comment{}

    if err := db.Preload("User").Where("user_id = ?", userID).Preload("Photo").Where("user_id = ?", userID).Find(&Comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to retrieve Comment",
        })
        return
    }

	var response []gin.H
    for _, comment := range Comment {
        response = append(response, gin.H{
			"id": comment.ID,
			"message": comment.Message,
			"photo_id": comment.PhotoID,
			"user_id": comment.UserID,
			"created_at": comment.CreatedAt,
			"updated_at": comment.UpdatedAt,
			"User": gin.H{
				"id": comment.User.ID,
				"email": comment.User.Email,
				"username": comment.User.UserName,
			},
			"Photo": gin.H{
				"id": comment.Photo.ID,
				"title": comment.Photo.Title,
				"caption": comment.Photo.Caption,
				"photo_url": comment.Photo.PhotoUrl,
				"user_id": comment.Photo.UserID,
			},
            
        })
    }

    c.JSON(http.StatusOK, response)
}

func UpdateComment(c *gin.Context) {
	db := database.GetDB()
	userData := c.MustGet("userData").(jwt.MapClaims)
	contentType := helpers.GetContentType(c)
	Comment := models.Comment{}

	commentId, _ := strconv.Atoi(c.Param("commentId"))
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&Comment)
	} else {
		c.ShouldBind(&Comment)
	}

	Comment.UserID = userID
	Comment.ID = uint(commentId)

	err := db.Model(&Comment).Where("id = ?", commentId).Updates(models.Comment{Message: Comment.Message}).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	response := db.Debug().Model(&Comment).Where("id = ?", commentId).First(&Comment)

	if response.Error != nil {
    c.JSON(http.StatusInternalServerError, gin.H{
        "error":   "Internal Server Error",
        "message": "Failed to retrieve comment",
    })
    return
}

	c.JSON(http.StatusOK, gin.H{
		"id": Comment.ID,
		"message": Comment.Message,
		"photo_id": Comment.PhotoID,
		"user_id": Comment.UserID,
		"updated_at": Comment.UpdatedAt,
	})
}

func DeleteComment(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    contentType := helpers.GetContentType(c)
    Comment := models.Comment{}

    commentId, err := strconv.Atoi(c.Param("commentId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Bad Request",
            "message": "Invalid comment ID",
        })
        return
    }

    id := uint(userData["id"].(float64))

    if contentType == appJSON {
        if err := c.ShouldBindJSON(&Comment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Bad Request",
                "message": err.Error(),
            })
            return
        }
    } else {
        if err := c.ShouldBind(&Comment); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{
                "error":   "Bad Request",
                "message": err.Error(),
            })
            return
        }
    }

    err = db.Where("id = ?", commentId).First(&Comment).Error
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "error":   "Not Found",
            "message": "Comment not found",
        })
        return
    }

    if uint(Comment.UserID) != id {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":   "Unauthorized",
            "message": "You are not authorized to delete this comment",
        })
        return
    }

    if err := db.Delete(&Comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Your comment has been successfully deleted",
    })
}
