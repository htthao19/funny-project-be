package authn

import (
	goctx "context"
	"fmt"
	"strconv"
	"strings"

	"github.com/beego/beego"
	"github.com/beego/beego/context"
	jwt "github.com/dgrijalva/jwt-go"

	"funny-project-be/infra/constant"
	"funny-project-be/infra/options"
)

// VerifyToken verifies the JWT.
func VerifyToken(opts options.Options) beego.FilterFunc {
	return func(ctx *context.Context) {
		if strings.HasPrefix(ctx.Input.URL(), "/funny-project/v1/rpc/auth/login") || strings.HasPrefix(ctx.Input.URL(), "/funny-project/v1/ws") {
			return
		}

		w := ctx.ResponseWriter
		claims, valid := IsValidJWT(opts, ctx.Request.Header.Get("Authorization"))
		if !valid {
			w.WriteHeader(401)
			w.Write([]byte("401 Unauthorized\n"))
			return
		}

		uid, err := strconv.Atoi(claims["sub"].(string))
		if err != nil {
			beego.Error("VerifyToken ", err)
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

func IsValidJWT(opts options.Options, token string) (jwt.MapClaims, bool) {
	// The Authorization header should come in this format: Bearer <jwt>
	s := strings.SplitN(token, " ", 2)
	if len(s) != 2 || s[0] != "Bearer" {
		return nil, false
	}

	t, err := jwt.Parse(s[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(opts.AccessTokenSecret), nil
	})
	if err != nil {
		beego.Error("IsValidJWT ", err)
		return nil, false
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok || !t.Valid {
		return nil, false
	}

	return claims, true
}
