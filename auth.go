package middleware

import (
	"net/http"

	"github.com/labstack/echo"
)

var (
	xUserID = http.CanonicalHeaderKey("X-User-Id")
)

// CurrentUserMiddleware 从 header 中获取当前用户信息
func CurrentUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get(xUserID)
		if userID != "" {
			c.Set(xUserID, userID)
		}

		return next(c)
	}
}

// AuthMiddleware 认证
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := c.Request().Header.Get(xUserID)
		if userID == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "请先登录")
		}

		return next(c)
	}
}
