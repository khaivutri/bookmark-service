package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler/v1/dto"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/internal/service"
	"github.com/rs/zerolog/log"
)

// ShortenURL defines the interface for handling URL shortening operations.
type ShortenURL interface {
	CreateShortenLink(ctx *gin.Context) 
	Redirect( ctx *gin.Context)
}

// shortenURL implements the ShortenURL interface.
type shortenURL struct {
	service service.ShortenURL
}


// NewShortenURL creates and returns a new instance of ShortenURL handler.
func NewShortenURL(service service.ShortenURL) ShortenURL {
	return &shortenURL{service: service}
}


// CreateShortenLink generates a shortened code for the provided URL and returns it in the response.
// @Summary      Create Shorten Link
// @Tags         ShortenURL
// @Accept       application/json
// @Produce      application/json
// @Param        request  body      dto.ShortenURLReq  true  "Shorten URL request"
// @Success      200  {object}  dto.ShortenURLRes  "Shorten URL response"
// @Failure      400  {object}  map[string]string  "Bad Request"
// @Failure      500  {object}  map[string]string  "Internal Server Error"
// @Router       /v1/links/shorten [post]
func (s *shortenURL) CreateShortenLink(ctx *gin.Context) {
	var req dto.ShortenURLReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	code, err := s.service.CreateCodeFromLink(ctx.Request.Context(), req.URL, req.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenURL.CreateShortenLink").Msg("failed to code from link")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	
	response := dto.ShortenURLRes{	Code: 		code,
									Message: 	"Shorten URL generated successfully!",
								}
	ctx.JSON(http.StatusOK, response)
}

// Redirect retrieves the original URL from the provided code and performs an HTTP redirect.
// @Summary      Redirect
// @Tags         ShortenURL
// @Accept       application/json
// @Produce      application/json
// @Param        code  path     string    	true  	"Code"
// @Success      302  {object}  map[string]string  	"Redirect"
// @Failure      400  {object}  map[string]string  	"Bad Request"
// @Failure      404  {object}  map[string]string  	"Not Found"
// @Failure      500  {object}  map[string]string  	"Internal Server Error"
// @Router       /v1/links/redirect/{code} [get]
func (s *shortenURL) Redirect( ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	url, err := s.service.GetLinkFromCode(ctx.Request.Context(), code)
	if err != nil {
		if errors.Is(err, repository.ErrCodeNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Code not found"})
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenURL.Redirect").Msg("failed to get link from code")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	
	ctx.Redirect(http.StatusFound, url)
}