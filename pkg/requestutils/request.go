package requestutils

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/khaivutri/bookmark-service/pkg/response"
)

var (
	InputValidator = validator.New(validator.WithRequiredStructEnabled())
)

func abortWithError(ctx *gin.Context, err error) error {
	ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldErrorResponse(err))
	return err
}

// BindInputFromResquest binds and validates input from various parts of the HTTP request (JSON, URI, query, headers).
func BindInputFromResquest[T any] (ctx *gin.Context) (*T, error) {
	reqInput := new(T)

	if ctx.Request.Method != http.MethodGet {
		if err := ctx.ShouldBindJSON(reqInput); err != nil && !errors.Is(err, io.EOF){
			return nil, abortWithError(ctx, err)
		}
	}	

	if err := ctx.ShouldBindUri(reqInput); err != nil {
		return nil, abortWithError(ctx, err)
	}

	if err := ctx.ShouldBindQuery(reqInput); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldErrorResponse(err))
		return nil, abortWithError(ctx, err)
	}

	if err := ctx.ShouldBindHeader(reqInput); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldErrorResponse(err))
		return nil, abortWithError(ctx, err)
	}

	if err := InputValidator.Struct(reqInput); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.InputFieldErrorResponse(err))
		return nil, abortWithError(ctx, err)
	}

	return reqInput, nil
}