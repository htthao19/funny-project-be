package controller

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	beegoctx "github.com/beego/beego/context"
	"github.com/stretchr/testify/assert"

	"funny-project-be/domain/entity"
	"funny-project-be/infra/constant"
)

type MockUserRepo struct{}

func (m *MockUserRepo) Get(ctx context.Context, uid uint) (*entity.User, error) {
	fmt.Println("Mock Get called")
	// Mock implementation
	return &entity.User{
		Email: "test@example.com",
	}, nil
}
func (m *MockUserRepo) GetOneByEmail(ctx context.Context, email string) (*entity.User, error) {
	return nil, nil
}
func (m *MockUserRepo) Add(ctx context.Context, users ...*entity.User) error {
	return nil
}
func (m *MockUserRepo) Remove(ctx context.Context, users ...*entity.User) error {
	return nil
}
func (m *MockUserRepo) Update(ctx context.Context, users ...*entity.User) error {
	return nil
}

func (c *UserController) ServeJSON() {}

func TestUserController_GetUser(t *testing.T) {
	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/user", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// Create a mock user repository
	mockRepo := &MockUserRepo{}

	// Create a UserController with the mock repository
	controller := &UserController{
		URepo: mockRepo,
		// Fill other dependencies here
	}

	// We need to set the context manually for testing
	controller.Ctx = &beegoctx.Context{
		Input:          beegoctx.NewInput(),
		Output:         beegoctx.NewOutput(),
		Request:        req,
		ResponseWriter: &beegoctx.Response{},
	}
	controller.Data = make(map[interface{}]interface{})
	controller.Ctx.Input.SetData(constant.ContextUID, uint(1))
	controller.Ctx.Input.SetData(constant.ContextCtx, context.Background())

	controller.GetUser()

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	expected := "test@example.com"
	assert.Equal(t, expected, controller.Data["json"].(*GetUserResponse).Email)
}
