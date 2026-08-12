package handler

import (
	"net/http"

	"github.com/d1zyy/monitor-pc/internal/buildinfo"

	"github.com/gin-gonic/gin"
)

func GetVersion(c *gin.Context) {
	versionInfo := buildinfo.Get()
	c.JSON(http.StatusOK, versionInfo)
}
