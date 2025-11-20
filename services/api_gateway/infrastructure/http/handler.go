package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/waste3d/ai-ops/services/api_gateway/application"
)

type Handler struct {
	useCase    *application.TicketUseCase
	jwtService *application.JWTService
	userClient *grpc_client.UserClient
	userReader application.UserReader
}

func NewHandler(useCase *application.TicketUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.GET("/tickets", h.GetAllTickets)
		api.GET("/tickets/:id", h.GetTicketByID)
	}
}

func (h *Handler) GetAllTickets(c *gin.Context) {
	tickets, err := h.useCase.GetAllTickets(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *Handler) GetTicketByID(c *gin.Context) {
	id := c.Param("id")
	ticket, err := h.useCase.GetTicketByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ticket)
}
