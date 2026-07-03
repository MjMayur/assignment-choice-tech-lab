package web

import (
	"net/http"

	projectContext "project/pkg/context"
	"project/pkg/log"
	httpUtils "project/pkg/utils"

	"github.com/gin-gonic/gin"
)

func (h CoreHandlerRegistry) HomeHandler(c *gin.Context) {
	reqID, _ := projectContext.GetRequestIDFromContext(c.Request.Context())
	log.Info("core>web>home: HomeHandler started", reqID)

	log.Info("core>web>login: HomeHandler completed", reqID)
	httpUtils.DataResponse(c, http.StatusOK, "Home successful.", nil)
}
