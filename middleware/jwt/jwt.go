package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/hespecial/go-gin-example/pkg/e"
	"github.com/hespecial/go-gin-example/pkg/util"
	"net/http"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
		} else {
			tokenString = tokenString[len(util.TokenType)+1:]
			tokenClaims, err := util.ParseToken(tokenString)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else {
				if tokenClaims.Issuer != util.Issuer {
					code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": e.GetMsg(code),
				"data":    make(map[string]string),
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
