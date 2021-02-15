package auth

import (
	"context"

	"github.com/go-redis/redis"
	"github.com/labstack/echo"
	"gorm.io/gorm"
)

type Handler struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func (h *Handler) Validator(key string, c echo.Context) (bool, error) {
	ctx := context.Background()

	if uid, err := h.RedisClient.Get(ctx, key).Result(); err != nil {
		return false, err
	} else {
		c.Set("uid", uid)
		return true, nil
	}
}
