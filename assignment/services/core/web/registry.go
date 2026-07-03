package web

import (
	"project/pkg/log"
	"project/pkg/server"
	coreService "project/services/core/service"

	"github.com/gin-gonic/gin"
)

type CoreHandlerRegistryOptions struct {
	Config                *server.Config
	AssignmentUserService coreService.AssignmentUserService
}

type CoreHandlerRegistry struct {
	Options CoreHandlerRegistryOptions
}

func NewCoreHandlerRegistry(options CoreHandlerRegistryOptions) *CoreHandlerRegistry {
	return &CoreHandlerRegistry{
		Options: options,
	}
}

func (h *CoreHandlerRegistry) StartServer() error {

	router, err := h.registerRoutes()
	if err != nil {
		log.Print(err)
	}

	log.Info("Server Started Successfully", "")
	router.Run(h.Options.Config.Port)
	return nil
}

func (h CoreHandlerRegistry) registerRoutes() (*gin.Engine, error) {

	router := gin.Default()

	coreRouter := router.Group("/tms-core")

	assignmentUserRouter := coreRouter.Group("/assignment-user")
	assignmentUserRouter.POST("/upload-csv", h.UploadCsv)
	assignmentUserRouter.POST("/", h.CreateAssignmentUserHandler)
	assignmentUserRouter.PATCH("/:id", h.PartialUpdateAssignmentUserHandler)
	assignmentUserRouter.DELETE("/:id", h.DeleteAssignmentUserHandler)
	assignmentUserRouter.GET("/:id", h.GetAssignmentUserbyIDHandler)
	assignmentUserRouter.GET("/", h.ListAssignmentUserHandler)
	//-----==-----==DO NOT ADD CODE BELOW THIS LINE------
	return router, nil
}
