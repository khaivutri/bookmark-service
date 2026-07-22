package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler/v1/dto"
	"github.com/khaivutri/bookmark-service/internal/repository"
	"github.com/khaivutri/bookmark-service/internal/service"
	"github.com/khaivutri/bookmark-service/pkg/response"
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


// CreateShortenLink generates a shortened code for the provided URL.
// @Summary      Create Shorten Link
// @Description  Creates a unique short code for the given URL and returns it.
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
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldErrorResponse(err))
		return
	}

	code, err := s.service.CreateCodeFromLink(ctx, req.URL, req.Exp)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.shortenURL.CreateShortenLink").Msg("failed to code from link")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	
	response := dto.ShortenURLRes{	Code: 		code,
									Message: 	"Shorten URL generated successfully!",
								}
	ctx.JSON(http.StatusOK, response)
}
// Redirect retrieves the original URL associated with the provided short code
// and redirects the client to that URL.
// @Summary      Redirect to original URL
// @Description  Resolves a short code and redirects the client to the original URL.
// @Tags         ShortenURL
// @Produce      plain
// @Param        code  path      string  true  "Short code"
// @Success      302   {string}  string  "Redirect to the original URL"
// @Failure      400   {object}  map[string]string  "Bad Request"
// @Failure      404   {object}  map[string]string  "Code not found"
// @Failure      500   {object}  map[string]string  "Internal Server Error"
// @Router       /v1/links/redirect/{code} [get]
func (s *shortenURL) Redirect( ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputErrResponse)
		return
	}

	url, err := s.service.GetLinkFromCode(ctx, code)
	if err != nil {
		if errors.Is(err, repository.ErrCodeNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Code not found"})
			return
		}

		log.Error().Err(err).Str("from", "handler.shortenURL.Redirect").Msg("failed to get link from code")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	
	ctx.Redirect(http.StatusFound, url)
}