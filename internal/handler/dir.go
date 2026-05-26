package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

type DirEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size"`
	ModTime string `json:"mod_time"`
}

func DirExplorer(c echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		path = "/"
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(err))
	}

	osEntries, err := os.ReadDir(abs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, errorResp(err))
	}

	entries := make([]DirEntry, 0, len(osEntries))
	for _, entry := range osEntries {
		info, _ := entry.Info()
		size := int64(0)
		modTime := ""
		if info != nil {
			size = info.Size()
			modTime = info.ModTime().Format(time.RFC3339)
		}
		entries = append(entries, DirEntry{
			Name:    entry.Name(),
			Path:    filepath.Join(abs, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    size,
			ModTime: modTime,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"current": abs,
		"entries": entries,
	})
}
