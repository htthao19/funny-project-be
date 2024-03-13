package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/beego/beego/logs"
	"github.com/beego/beego/validation"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"

	"funny-project-be/domain/entity"
	"funny-project-be/domain/repo"
	"funny-project-be/infra/constant"
	"funny-project-be/infra/options"
	"funny-project-be/infra/status"
)

// UserController exposes apis of User resource.
type UserController struct {
	BaseController

	URepo repo.UserRepo

	Opts options.Options
}

// User info.
type User struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

// NewUserFromEntity creates User from entity.
func NewUserFromEntity(e *entity.User) *User {
	if e == nil {
		return nil
	}

	return &User{
		ID:        e.ID,
		Name:      e.Name,
		Email:     e.Email,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// LoginRequest is a struct contains a login request.
type LoginRequest struct {
	Code        string `json:"code" valid:"Required"`
	RedirectURL string `json:"redirectURL" valid:"Required"`
}

// LoginResponse is a struct contains a login reponse.
type LoginResponse struct {
	Response
	Token  string `json:"token,omitempty"`
	Name   string `json:"name,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

// GoogleUser represents for google user infomation.
type GoogleUser struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

// JWTClaim is the JWT custom claims.
type JWTClaim struct {
	Email string `json:"email,omitempty"`
	jwt.StandardClaims
}

// Login API.
func (c *UserController) Login() {
	var req LoginRequest
	var resp LoginResponse
	resp.Code = status.OK

	defer func() {
		c.Ctx.Output.SetStatus(resp.Code / 1000)
		c.Data["json"] = &resp
		c.ServeJSON()
	}()

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil {
		resp.Code = status.BadRequest
		resp.SetError(err)
		return
	}

	var validator validation.Validation
	valid, err := validator.Valid(&req)
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		logs.Error("Login ", err)
		return
	}
	if !valid {
		resp.Code = status.BadRequest
		resp.SetValidationErrors(validator.Errors)
		return
	}

	conf := &oauth2.Config{
		ClientID:     c.Opts.GAuthClientID,
		ClientSecret: c.Opts.GAuthClientSecret,
		RedirectURL:  req.RedirectURL,
		Endpoint:     google.Endpoint,
	}
	fmt.Printf("%+v\n", conf)

	// exchange auth_code to token including refresh_token
	token, err := conf.Exchange(context.Background(), req.Code)
	if err != nil {
		resp.Code = status.Unauthorized
		resp.Message = fmt.Sprintf(`exchange google token failed %s`, err.Error())
		return
	}

	client := conf.Client(context.Background(), token)
	response, err := client.Get(c.Opts.GAuthProfileURL)
	if err != nil {
		resp.Code = status.Unauthorized
		resp.Message = fmt.Sprintf(`get google user failed %s`, err.Error())
		return
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		resp.Message = fmt.Sprintf(`get google user failed %s`, err.Error())
		resp.Code = status.Unauthorized
		return
	}
	fmt.Println(string(contents))

	gUser := GoogleUser{}
	if err = json.Unmarshal(contents, &gUser); err != nil {
		resp.Message = fmt.Sprintf(`parse google user failed %s`, err.Error())
		resp.Code = status.Unauthorized
		return
	}

	ctx := context.Background()
	user, err := c.URepo.GetOneByEmail(ctx, gUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &entity.User{
				Email: gUser.Email,
				Name:  gUser.Name,
			}
			if err := c.URepo.Add(ctx, user); err != nil {
				resp.Code = status.InternalServerError
				resp.SetError(err)
				logs.Error("Login ", err)
				return
			}
			user, err = c.URepo.GetOneByEmail(ctx, gUser.Email)
			if err != nil {
				resp.Code = status.InternalServerError
				resp.SetError(err)
				logs.Error("Login ", err)
				return
			}
		} else {
			resp.Code = status.InternalServerError
			resp.SetError(err)
			logs.Error("Login ", err)
			return
		}
	}

	expiredAt := time.Now().Unix() + int64(c.Opts.AccessTokenExpiresIn.Seconds())
	claims := JWTClaim{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expiredAt,
			Subject:   strconv.Itoa(int(user.ID)),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	signedStr, err := accessToken.SignedString([]byte(c.Opts.AccessTokenSecret))
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		return
	}

	resp.Token = signedStr
	resp.Name = gUser.Name
	resp.Avatar = gUser.Picture
}

// GetUserResponse is a response of GetUser API.
type GetUserResponse struct {
	Response
	Email string `json:"email,omitempty"`
}

// GetUser API.
func (c *UserController) GetUser() {
	var resp GetUserResponse
	resp.Code = status.OK

	defer func() {
		c.Ctx.Output.SetStatus(resp.Code / 1000)
		c.Data["json"] = &resp
		c.ServeJSON()
	}()

	uid := c.Ctx.Input.GetData(constant.ContextUID).(uint)
	ctx := c.Ctx.Input.GetData(constant.ContextCtx).(context.Context)
	user, err := c.URepo.Get(ctx, uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = status.NotFound
			resp.Message = `user not found`
			return
		}
		resp.Code = status.InternalServerError
		resp.SetError(err)
		logs.Error("GetUser ", err)
		return
	}
	resp.Email = user.Email
}
