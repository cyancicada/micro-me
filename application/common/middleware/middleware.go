package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"micro-me/application/common/baseerror"
	"micro-me/application/common/baseresponse"
)

var (
	DefaultField                           = "Authorization"
	AccessTokenValidErr                    = baseerror.NewBaseError("AccessToken 验证失败")
	AccessTokenValidationErrorExpiredErr   = baseerror.NewBaseError("AccessToken过期")
	AccessTokenValidationErrorMalformedErr = baseerror.NewBaseError("AccessToken格式错误")
)

const (
	UserSignedKey = "vector_"
)

func ValidAccessToken(context *gin.Context) {
	authorization := context.GetHeader(DefaultField)
	token, err := jwt.Parse(authorization, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(UserSignedKey), nil
	})
	if err != nil {
		if err, ok := err.(*jwt.ValidationError); ok {
			if err.Errors&jwt.ValidationErrorMalformed != 0 {
				baseresponse.HttpResponse(context, nil, AccessTokenValidationErrorMalformedErr)
				context.Abort()
				return
			}
			if err.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				baseresponse.HttpResponse(context, nil, AccessTokenValidationErrorExpiredErr)
				context.Abort()
				return
			}
		}
		baseresponse.HttpResponse(context, nil, AccessTokenValidErr)
		context.Abort()
		return
	}
	if token != nil && token.Valid {
		context.Next()
		return
	}
	baseresponse.HttpResponse(context, nil, AccessTokenValidErr)
	context.Abort()
	return

}
