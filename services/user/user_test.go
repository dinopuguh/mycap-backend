package user_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/routes"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/stretchr/testify/assert"
)

var (
	createdUser user.User
)

func TestNew(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	app := routes.New()

	type args struct {
		data        user.RegisterUser
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid register", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@email.com",
				Username: "dinopuguh16",
				Password: "12345678",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Email already exist.", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@email.com",
				Username: "dinopuguh16",
				Password: "12345678",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Username already exist.", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@gmail.com",
				Username: "dinopuguh16",
				Password: "12345678",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@gmail.com",
				Username: "dinopuguh16",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				rb := user.ResponseAuth{}
				json.Unmarshal(resBody, &rb)

				userJson, _ := json.Marshal(rb.User)
				json.Unmarshal(userJson, &createdUser)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	router := routes.New()
	t.Run("Get all users", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)

		res, _ := router.Test(req, -1)
		defer res.Body.Close()
		resBody, _ := ioutil.ReadAll(res.Body)

		assert.Equalf(t, http.StatusOK, res.StatusCode, string(resBody))
	})
}

func TestLogin(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	r := routes.New()

	type args struct {
		data        user.LoginUser
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid login", args{
			data: user.LoginUser{
				Email:    "dino@email.com",
				Password: "12345678",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"User not found", args{
			data: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "12345678",
			},
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Password incorrect", args{
			data: user.LoginUser{
				Email:    "dino@email.com",
				Password: "123456789",
			},
			statusCode:  http.StatusUnauthorized,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: user.LoginUser{
				Email:    "dino@email.com",
				Password: "12345678",
			},
			statusCode: http.StatusBadRequest,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}
