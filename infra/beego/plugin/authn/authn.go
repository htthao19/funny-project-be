package authn

import (
	goctx "context"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego"
	"github.com/beego/beego/context"
	"github.com/beego/beego/logs"
	jwt "github.com/dgrijalva/jwt-go"

	"funny-project-be/infra/constant"
	"funny-project-be/infra/options"
)

// VerifyToken verifies the JWT.
func VerifyToken(opts options.Options) beego.FilterFunc {
	return func(ctx *context.Context) {
		if strings.HasPrefix(ctx.Input.URL(), "/funny-project/v1/rpc/auth/login") {
			return
		}

		w := ctx.ResponseWriter
		// The Authorization header should come in this format: Bearer <jwt>
		s := strings.SplitN(ctx.Request.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 || s[0] != "Bearer" {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}

		tokenString := s[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(opts.AccessTokenSecret), nil
		})
		if err != nil {
			logs.Error("VerifyToken ", err)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}

		uid, err := strconv.Atoi(claims["sub"].(string))
		if err != nil {
			logs.Error("VerifyToken ", err)
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}
		ctx.Input.SetData(constant.ContextUID, uint(uid))
		ctx.Input.SetData(constant.ContextEmail, claims["email"])

		customctx := goctx.Background()
		customctx = goctx.WithValue(customctx, constant.ContextUID, uint(uid))
		customctx = goctx.WithValue(customctx, constant.ContextEmail, claims["email"])
		ctx.Input.SetData(constant.ContextCtx, customctx)
	}
}
