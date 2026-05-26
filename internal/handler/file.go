package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func FileServe(c echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		return c.JSON(http.StatusBadRequest, errorResp(os.ErrInvalid))
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(err))
	}

	info, err := os.Stat(abs)
	if err != nil {
		return c.JSON(http.StatusNotFound, errorResp(err))
	}
	if info.IsDir() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "path is a directory"})
	}

	return c.File(abs)
}
