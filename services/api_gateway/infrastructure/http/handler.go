package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/waste3d/ai-ops/services/api_gateway/application"
	grpc_client "github.com/waste3d/ai-ops/services/api_gateway/infrastructure/gprc"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	useCase    *application.TicketUseCase
	jwtService *application.JWTService
	userClient *grpc_client.UserClient
	userReader application.UserReader
}

func NewHandler(useCase *application.TicketUseCase, jwtService *application.JWTService, userClient *grpc_client.UserClient, userReader application.UserReader) *Handler {
	return &Handler{useCase: useCase, jwtService: jwtService, userClient: userClient, userReader: userReader}
}

func (h *Handler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format, expected 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		claims, err := h.jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		authorized := api.Group("/")
		authorized.Use(h.AuthMiddleware())
		{
			// Этот эндпоинт теперь защищен
			authorized.GET("/tickets", h.GetAllTickets)
			authorized.GET("/tickets/:id", h.GetTicketByID)
		}
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

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=4"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" binding:"required,min=8"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password and confirm password do not match"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to process password"})
		return
	}

	user, err := h.userClient.Register(c.Request.Context(), req.Username, string(passwordHash))
	if err != nil {
		log.Printf("Error registering user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.Id, "username": user.Username})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userReader.GetUserByUsername(c.Request.Context(), req.Username)
	if err != nil {
		log.Printf("Error getting user by username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		log.Printf("Error comparing password: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	token, err := h.jwtService.GenerateToken(user.Id)
	if err != nil {
		log.Printf("Error generating token: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
