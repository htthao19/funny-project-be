package controller

import (
	"container/list"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/beego/beego"
	"github.com/beego/beego/validation"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"

	"funny-project-be/domain/entity"
	"funny-project-be/domain/repo"
	"funny-project-be/infra/beego/plugin/authn"
	"funny-project-be/infra/constant"
	"funny-project-be/infra/options"
	"funny-project-be/infra/status"
)

// VideoController exposes apis of Video resource.
type VideoController struct {
	BaseController

	VRepo repo.VideoRepo
	URepo repo.UserRepo

	Opts options.Options
}

// Video info.
type Video struct {
	ID          uint      `json:"id,omitempty"`
	URL         string    `json:"url,omitempty"`
	SharedBy    string    `json:"sharedBy,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// NewVideoFromEntity creates Video from entity.
func NewVideoFromEntity(e *entity.Video) *Video {
	if e == nil {
		return nil
	}

	return &Video{
		ID:          e.ID,
		URL:         e.URL,
		SharedBy:    e.SharedBy,
		Description: e.Description,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}

type Subscriber struct {
	UID  uint
	IP   string
	Conn *websocket.Conn // Only for WebSocket users; otherwise nil.
}

var (
	// Channel for new join users.
	subscribe   = make(chan Subscriber, 100)
	subscribers = list.New()
)

// GetVideoRequest represents a request for get video.
type GetVideoRequest struct {
	ID uint
}

// GetVideoResponse is a response of GetVideo API.
type GetVideoResponse struct {
	Response
	*Video
}

// GetVideo API.
func (c *VideoController) GetVideo() {
	var req GetVideoRequest
	var resp GetVideoResponse
	resp.Code = status.OK

	defer func() {
		c.Ctx.Output.SetStatus(resp.Code / 1000)
		c.Data["json"] = &resp
		c.ServeJSON()
	}()

	id, err := strconv.ParseUint(c.Ctx.Input.Param(":id"), 10, 64)
	if err != nil {
		resp.Code = status.BadRequest
		resp.Message = "id is invalid"
		return
	}
	req.ID = uint(id)
	ctx := c.Ctx.Input.GetData(constant.ContextCtx).(context.Context)
	video, err := c.VRepo.Get(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			resp.Code = status.NotFound
			resp.Message = `video not found`
			return
		}
		resp.Code = status.InternalServerError
		beego.Error("GetVideo ", err)
		return
	}

	resp.Video = NewVideoFromEntity(video)
}

// CreateVideoRequest is a request of CreateVideo.
type CreateVideoRequest struct {
	URL         string `json:"url,omitempty" valid:"Required"`
	Description string `json:"description,omitempty"`
}

// CreateVideoResponse is a struct for return new created video.
type CreateVideoResponse struct {
	Response
	Video *Video `json:"video,omitempty"`
}

// CreateVideo API.
func (c *VideoController) CreateVideo() {
	var req CreateVideoRequest
	var resp CreateVideoResponse
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
		beego.Error("CreateVideo ", err)
		return
	}
	if !valid {
		resp.Code = status.BadRequest
		resp.SetValidationErrors(validator.Errors)
		return
	}

	video := &entity.Video{
		URL:         req.URL,
		Description: req.Description,
	}
	uid := c.Ctx.Input.GetData(constant.ContextUID).(uint)
	ctx := c.Ctx.Input.GetData(constant.ContextCtx).(context.Context)

	user, err := c.URepo.Get(ctx, uid)
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		beego.Error("GetUser ", err)
		return
	}
	video.SharedBy = user.Email

	if err := c.VRepo.Add(ctx, video); err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		beego.Error("CreateVideo ", err)
		return
	}

	c.broadcastWebSocket(NewVideoFromEntity(video), uid)

	resp.Code = status.Created
	resp.Video = NewVideoFromEntity(video)
}

// ListVideosRequest represents a request for listing videos.
type ListVideosRequest struct {
	Page  int    `form:"page" valid:"Min(1)"`
	Limit int    `form:"limit" valid:"Range(1, 200)"`
	Sort  string `form:"sort"`
}

// ListVideosResponse is the response of ListVideos.
type ListVideosResponse struct {
	RangeResponse
	Items []*Video `json:"_items"`
}

// ListVideos API.
func (c *VideoController) ListVideos() {
	var req ListVideosRequest
	var resp ListVideosResponse
	resp.Code = status.OK
	resp.Items = []*Video{}

	defer func() {
		c.Ctx.Output.SetStatus(resp.Code / 1000)
		c.Data["json"] = &resp
		c.ServeJSON()
	}()

	if err := c.ParseForm(&req); err != nil {
		resp.Code = status.BadRequest
		resp.SetError(err)
		return
	}

	var validator validation.Validation
	valid, err := validator.Valid(&req)
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		beego.Error("ListVideos ", err)
		return
	}
	if !valid {
		resp.Code = status.BadRequest
		resp.SetValidationErrors(validator.Errors)
		return
	}

	ctx := c.Ctx.Input.GetData(constant.ContextCtx).(context.Context)
	videos, err := c.VRepo.GetRangeByQuery(ctx, req.Sort, req.Limit, req.Page)
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		beego.Error("ListVideos ", err)
		return
	}
	total, err := c.VRepo.Count(ctx)
	if err != nil {
		resp.Code = status.InternalServerError
		resp.SetError(err)
		beego.Error("ListVideos ", err)
		return
	}

	resp.Page = req.Page
	resp.Limit = req.Limit
	resp.Total = total
	for _, u := range videos {
		if video := NewVideoFromEntity(u); u != nil {
			resp.Items = append(resp.Items, video)
		}
	}
}

// Join method handles WebSocket requests.
func (c *VideoController) JoinWebSocket() {
	ip := c.Ctx.Input.IP()
	// Upgrade from http request to WebSocket.
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// Allow all origins for WebSocket connections
			return true
		},
	}
	ws, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	_, tokenBytes, err := ws.ReadMessage()
	if err != nil {
		ws.Close()
		return
	}
	claims, valid := authn.IsValidJWT(c.Opts, string(tokenBytes))
	if !valid {
		ws.Close()
		return
	}
	uid, err := strconv.Atoi(claims["sub"].(string))
	if err != nil {
		ws.Close()
		return
	}
	subscribe <- Subscriber{UID: uint(uid), IP: ip, Conn: ws}
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func (c *VideoController) broadcastWebSocket(video *Video, uid uint) {
	data, err := json.Marshal(video)
	if err != nil {
		beego.Error("Fail to marshal video:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		if uid == sub.Value.(Subscriber).UID {
			continue
		}
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				beego.Info("User disconnected")
			}
		}
	}
}

// This function handles all incoming chan messages.
func handleScription() {
	for {
		sub := <-subscribe
		subscribers.PushBack(sub)
		fmt.Printf("%v+\n", subscribers)
	}
}

func init() {
	go handleScription()
}
