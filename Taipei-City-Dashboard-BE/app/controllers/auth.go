package controllers

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"TaipeiCityDashboardBE/app/models"
	"TaipeiCityDashboardBE/app/util"
	"TaipeiCityDashboardBE/global"
	"TaipeiCityDashboardBE/logs"

	"github.com/gin-gonic/gin"
)

const (
	emailRegex = "^\\w+((-\\w+)|(\\.\\w+))*\\@[A-Za-z0-9]+((\\.|-)[A-Za-z0-9]+)*\\.[A-Za-z]+$"
)

func Login(c *gin.Context) {
	var user models.AuthUser

	const authPrefix = "Basic "

	credentials, err := util.GetAuthFromRequest(c, authPrefix)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	email, password, err := util.DecodeCredentials(credentials)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// check parameters
	emailRegexp, err := regexp.MatchString(emailRegex, email)
	if err != nil || !emailRegexp {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Errorf("invalid email format: %v", err)})
		return
	}

	if password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password is required"})
		return
	}

	// Email's unique constraint is case-sensitive while login is not. Verify
	// candidates in ID order, preserving the old first matching account behavior.
	var candidates []models.AuthUser
	if err := models.DBManager.
		Where("LOWER(email) = LOWER(?)", email).
		Order("id ASC").
		Find(&candidates).Error; err != nil {
		logs.FError("Login failed: unexpected database error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unexpected database error"})
		return
	}
	validPassword, needsUpgrade := false, false
	for _, candidate := range candidates {
		if candidate.Password == nil {
			continue
		}
		if valid, upgrade := util.VerifyPassword(password, *candidate.Password); valid {
			user = candidate
			validPassword, needsUpgrade = true, upgrade
			break
		}
	}
	if !validPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect username or password"})
		return
	}

	// check user is active
	if user.IsActive == nil || !*user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "User not activated"})
		return
	}

	permissions, err := models.GetUserPermission(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// generate JWT token
	user.LoginAt = time.Now()
	token, err := util.GenerateJWT(user.LoginAt.Add(global.TokenExpirationDuration), "Email", user.ID, *user.IsAdmin, permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Upgrade existing SHA-256 credentials together with the login timestamp.
	// Keeping the replacement on user prevents Save from restoring the old hash.
	if needsUpgrade {
		replacement, err := util.HashPassword(password)
		if err != nil {
			logs.FError("Failed to upgrade password hash: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected password hashing error"})
			return
		}
		user.Password = &replacement
	}

	// update last login time and any password upgrade before returning the token
	if err := models.DBManager.Save(&user).Error; err != nil {
		logs.FError("Failed to update login time: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected database error"})
		return
	}

	// return JWT token
	c.JSON(http.StatusOK, gin.H{
		"user":       user,
		"permission": permissions,
		"token":      token,
	})
}
