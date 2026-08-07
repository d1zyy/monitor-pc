package handler

import (
	"monitor-pc/internal/buildinfo"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	versionInfo := buildinfo.Get()
	c.JSON(http.StatusOK, versionInfo)
}
