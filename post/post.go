package post

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"example.com/social-gin/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Post represents user post
type Post struct {
	ID        uint         `gorm:"primarykey" json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"update_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"-"`
	UserID    int          `json:"user_id"`
	User      user.User    `json:"-"`
	Content   string       `json:"content"`
	Likes     int          `json:"likes"`
}

// Handler handles user requests
type Handler struct {
	DB *gorm.DB
}

// AddPost handle add post request
func (h *Handler) AddPost(c *gin.Context) {
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
	post := Post{}

	if err := c.Bind(&post); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	post.UserID = uid

	if result := h.DB.Create(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, post)
}

// ListPost handle list post request
func (h *Handler) ListPost(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	posts := []Post{}
	if result := h.DB.Where("user_id = ?", uid).Find(&posts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, posts)
}

// GetPost handle get post request
func (h *Handler) GetPost(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	post := Post{}
	if result := h.DB.Where("user_id = ? and id = ?", uid, pid).First(&post); result.Error != nil {
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
	c.JSON(http.StatusOK, post)
}

// UpdatePost handle update post request
func (h *Handler) UpdatePost(c *gin.Context) {
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
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	post := Post{}
	if result := h.DB.Where("user_id = ? and id = ?", uid, pid).First(&post); result.Error != nil {
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

	updatePost := Post{}
	if err := c.Bind(&updatePost); err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// update fields
	if post.Content != "" {
		post.Content = updatePost.Content
	}
	if post.Likes != 0 {
		post.Likes = updatePost.Likes
	}

	if result := h.DB.Save(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, post)
}

// DeletePost handle delete post request
func (h *Handler) DeletePost(c *gin.Context) {
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
	pid, err := strconv.Atoi(c.Param("pid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	post := Post{}
	if result := h.DB.Where("user_id = ? and id = ?", uid, pid).First(&post); result.Error != nil {
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

	if result := h.DB.Delete(&post); result.Error != nil {
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, post)
}
