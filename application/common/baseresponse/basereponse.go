package baseresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"micro-me/application/common/baseerror"
	"micro-me/application/common/i18n"
)

func ParamError(ctx *gin.Context, err interface{}) {
	validErr, ok := err.(validator.ValidationErrors)
	if ok {
		errMap := map[string]string{}
		for _, ve := range validErr {
			key := ve.FieldNamespace + "." + ve.Tag
			errMap[key] = i18n.ZhMessage[key]
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"message": errMap})
		return
	}
	ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.ErrParam})
	return
}

func HttpResponse(ctx *gin.Context, res, err interface{}) {
	baeError, ok := err.(*baseerror.BaseError)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"message": baeError.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": i18n.ErrServer})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": res})
	return
}
