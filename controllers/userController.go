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

var (
	appJSON = "application/json"
)

func UserRegister(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	err := db.Debug().Create(&User).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"age": User.Age,
		"email": User.Email,
		"id": User.ID,
		"username": User.UserName,
	})	
}

func UserLogin(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	_, _ = db, contentType
	User := models.User{}
	password := ""

	if contentType == appJSON {
		c.ShouldBindJSON(&User)
	} else {
		c.ShouldBind(&User)
	}

	password = User.Password

	err := db.Debug().Where("email = ?", User.Email).Take(&User).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.ComparePass([]byte(User.Password), []byte(password))

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
			"message": "invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(User.ID, User.UserName, User.Email, User.Age)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func UpdateUser(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    contentType := helpers.GetContentType(c)
    User := models.User{}

    userId, _ := strconv.Atoi(c.Param("userId"))
	userID := uint(userData["id"].(float64))

    if contentType == appJSON {
        c.ShouldBindJSON(&User)
    } else {
        c.ShouldBind(&User)
    }

    User.ID = uint(userID)
    User.Age = int(userData["age"].(float64))

    err := db.Model(&User).Where("id = ?", userId).Updates(models.User{UserName: User.UserName, Email: User.Email}).Error
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Bad Request",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":         User.ID,
        "email":      User.Email,
        "username":   User.UserName,
        "age":        User.Age,
        "update_at":  User.UpdatedAt,
    })
}


func DeleteUser(c *gin.Context) {
    db := database.GetDB()
    userData := c.MustGet("userData").(jwt.MapClaims)
    userID := uint(userData["id"].(float64))
	
    requestedUserID, err := strconv.Atoi(c.Param("userId"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Bad Request",
            "message": "Invalid user ID",
        })
        return
    }
    if uint(requestedUserID) != userID {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error":   "Unauthorized",
            "message": "You are not authorized to delete this user",
        })
        return
    }

    
    tx := db.Begin()

    
    if err := tx.Exec("DELETE FROM social_media WHERE user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to delete associated social_media",
        })
        return
    }

    if err := tx.Exec("DELETE FROM comments WHERE user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to delete associated comments",
        })
        return
    }

    
    if err := tx.Exec("DELETE FROM photos WHERE user_id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to delete associated photos",
        })
        return
    }

    
    if err := tx.Exec("DELETE FROM users WHERE id = ?", userID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Internal Server Error",
            "message": "Failed to delete user",
        })
        return
    }

    
    tx.Commit()

    c.JSON(http.StatusOK, gin.H{
        "message": "Your account has been successfully deleted",
    })
}


