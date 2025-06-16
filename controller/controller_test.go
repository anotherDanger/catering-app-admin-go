package controller

import (
	"bytes"
	"catering-admin-go/domain"
	"catering-admin-go/service/mocks"
	"catering-admin-go/web"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(svc *mocks.Service)
		expectedErr    bool
		expectedResult *web.AdminResponse
	}{
		{
			name: "Success",

			setupMock: func(svc *mocks.Service) {
				rows := &web.AdminResponse{
					Username:    "ADMIN",
					AccessToken: "valid",
				}
				svc.On("Login", mock.Anything, mock.Anything).Return(rows, nil)
			},
			expectedErr: false,
			expectedResult: &web.AdminResponse{
				Username:    "ADMIN",
				AccessToken: "valid",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			svc := mocks.NewService(t)
			tt.setupMock(svc)
			ctrl := NewControllerImpl(svc)

			app.Post("/v1/login", ctrl.Login)

			reqBody := &domain.Admin{
				Username: "admin",
				Password: "admin",
			}

			jsonBytes, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(jsonBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}
			if tt.expectedErr {
				assert.NotEqual(t, fiber.StatusOK, resp.StatusCode)
			} else {
				assert.Equal(t, fiber.StatusOK, resp.StatusCode)

				var resBody web.AdminResponse
				err := json.NewDecoder(resp.Body).Decode(&resBody)
				assert.NoError(t, err)

				assert.Equal(t, tt.expectedResult.Username, resBody.Username)
				assert.Equal(t, tt.expectedResult.AccessToken, resBody.AccessToken)
			}
		})
	}
}

func TestDeleteOrder(t *testing.T) {
	id := "1"
	tests := []struct {
		name        string
		setupMock   func(svc *mocks.Service)
		expectedErr bool
	}{
		{
			name: "Success",
			setupMock: func(svc *mocks.Service) {
				svc.On("DeleteOrder", mock.Anything, mock.Anything).Return(nil)
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			svc := mocks.NewService(t)
			tt.setupMock(svc)
			ctrl := NewControllerImpl(svc)

			app.Delete("/api/v1/orders/:id", ctrl.DeleteOrder)
			req := httptest.NewRequest(fiber.MethodDelete, "/api/v1/orders/"+id, nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		})
	}
}
