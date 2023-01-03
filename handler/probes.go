package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// LivenessProbe implements a Health Check probe, responsible for reporting on
// application health.  It is used by Kubernetes in a Liveness Probe.
func LivenessProbe(c echo.Context) error {
	// TODO: Check DB connection is healthy
	return c.String(http.StatusOK, "I'm alive!")
}
