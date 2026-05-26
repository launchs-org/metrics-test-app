package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"metrics-sample-app/internal/stress"
)

type StressHandler struct {
	state *stress.State
}

func NewStressHandler(state *stress.State) *StressHandler {
	return &StressHandler{state: state}
}

func (h *StressHandler) Status(c echo.Context) error {
	return c.JSON(http.StatusOK, h.state.GetStatus())
}

type cpuStartReq struct {
	Percent float64 `json:"percent"`
}

func (h *StressHandler) CPUStart(c echo.Context) error {
	var req cpuStartReq
	if err := c.Bind(&req); err != nil || req.Percent <= 0 || req.Percent > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "percent must be 1-100"})
	}
	h.state.StartCPU(req.Percent)
	return c.JSON(http.StatusOK, map[string]string{"status": "started"})
}

func (h *StressHandler) CPUStop(c echo.Context) error {
	h.state.StopCPU()
	return c.JSON(http.StatusOK, map[string]string{"status": "stopped"})
}

type memStartReq struct {
	MB int64 `json:"mb"`
}

func (h *StressHandler) MemStart(c echo.Context) error {
	var req memStartReq
	if err := c.Bind(&req); err != nil || req.MB <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "mb must be > 0"})
	}
	h.state.StartMem(req.MB)
	return c.JSON(http.StatusOK, map[string]string{"status": "started"})
}

func (h *StressHandler) MemStop(c echo.Context) error {
	h.state.StopMem()
	return c.JSON(http.StatusOK, map[string]string{"status": "stopped"})
}
