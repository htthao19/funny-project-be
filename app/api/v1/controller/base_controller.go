package controller

import (
	"github.com/beego/beego"
	"github.com/beego/beego/validation"
)

// BaseController creates BaseController.
type BaseController struct {
	beego.Controller
}

// Prepare runs before controller.
func (c *BaseController) Prepare() {
}

// Response represents a base response for controller.
type Response struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// RangeResponse represents a base response for controller.
type RangeResponse struct {
	Response
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// SetError sets error to response.
func (r *Response) SetError(err error) {
	r.Message = err.Error()
}

// SetValidationErrors sets errors validation to responser
func (r *Response) SetValidationErrors(errs []*validation.Error) {
	for _, err := range errs {
		r.Message += err.Key + ":" + err.Message + ";"
	}
}
