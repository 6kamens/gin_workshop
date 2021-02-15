package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/social-gin/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// User represents user data
type User struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"update_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	Username  string       `gorm:"uniquekey" json:"username"`
	Password  string       `json:"password"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
}

// Handler represents handler of user data
type Handler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

// Hello handles hello request
func (h *Handler) Hello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}

// ListTable handles list table request
func (h *Handler) ListTable(c *gin.Context) {

	// custom logger
	logger := logger.Extract(c)
	uid := c.GetString("uid")
	logger.Info("listing table", zap.String("uid", uid))

	rows, err := h.DB.Raw("sp_tables").Rows()
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tables := []string{}
	var tableQualifier sql.NullString
	var tableOwner sql.NullString
	var tableName sql.NullString
	var tableType sql.NullString
	var remarks sql.NullString
	for rows.Next() {

		if err := rows.Scan(&tableQualifier, &tableOwner, &tableName, &tableType, &remarks); err != nil {
			c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		tables = append(tables, tableName.String)
	}
	c.JSON(http.StatusOK, tables)
}

// AddUser handle add user request
func (h *Handler) AddUser(c *gin.Context) {
	user := User{}
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	if result := h.DB.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// ListUser handle list user request
func (h *Handler) ListUser(c *gin.Context) {
	users := []User{}
	if result := h.DB.Find(&users); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser handle list user request
func (h *Handler) GetUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handle update user request
func (h *Handler) UpdateUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	loginUid := c.MustGet("uid")

	// c.Get("uid").(string)

	if loginUid != c.Param("uid") {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized user",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}

	updateUser := User{}
	if err := c.Bind(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// update fields
	if updateUser.Username != "" {
		user.Username = updateUser.Username
	}
	if updateUser.Password != "" {
		user.Password = updateUser.Password
	}
	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	if result := h.DB.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handle delete user request
func (h *Handler) DeleteUser(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	loginUid := c.MustGet("uid")

	if loginUid != c.Param("uid") {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized user",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	user := User{}
	if result := h.DB.First(&user, uid); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}

	if result := h.DB.Delete(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// LogIn handle login request
func (h *Handler) LogIn(c *gin.Context) {

	username := c.Request.FormValue("u")
	password := c.Request.FormValue("p")

	l := logger.Extract(c)

	l.Info(fmt.Sprint("login username=", username, " password=", password))

	user := User{}

	if result := h.DB.Where("Username = ?", username).Limit(1).Find(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid username or password",
		})
		return
	}

	if user.Password != password {
		c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid username or password",
		})
		return
	}

	token := uuid.New().String()

	l.Info(fmt.Sprint("new token=", token))

	if _, err := h.RedisClient.Set(token, user.ID, time.Minute*10).Result(); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

// Authorize authenticate using authorization header
func (h *Handler) Authorize(c *gin.Context) {

	auth := c.GetHeader("Authorization")
	prefix := "Bearer "

	if !strings.HasPrefix(auth, prefix) {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "no authorization token found in the header",
		})
		c.Abort()
		return
	}

	token := auth[len(prefix):]

	// validate token

	var uid string
	var err error
	if uid, err = h.RedisClient.Get(token).Result(); err != nil {
		// can't connect to redis
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "invalid token",
			})
		} else {
			c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": "can't connect to redis server",
			})
		}
		c.Abort()
		return
	}

	// set user id of the authenticated user to context
	c.Set("uid", uid)
}
