package handler

import (
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/labstack/echo/v4"
)

type EnvVar struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func Env(c echo.Context) error {
	envs := os.Environ()
	sort.Strings(envs)

	vars := make([]EnvVar, 0, len(envs))
	for _, entry := range envs {
		key, val, _ := strings.Cut(entry, "=")
		vars = append(vars, EnvVar{Key: key, Value: val})
	}
	return c.JSON(http.StatusOK, vars)
}

func Workdir(c echo.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, errorResp(err))
	}
	return c.JSON(http.StatusOK, map[string]string{"path": wd})
}
