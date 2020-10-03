package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/routes"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/stretchr/testify/assert"
)

var (
	createdUser *user.User
	updatedUser *user.User
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
		willUpdate  bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid register 1", args{
			data: user.RegisterUser{
				Name:     "Dino Puguh",
				Email:    "dinopuguh@mycap.com",
				Username: "dinopuguh",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Valid register 2", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@email.com",
				Username: "dinopuguh16",
				Password: "12345678",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Valid register 3", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dinopuguh@email.com",
				Username: "dino16",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
			willUpdate:  true,
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

				if tt.args.willUpdate {
					json.Unmarshal(userJson, &updatedUser)
				} else {
					json.Unmarshal(userJson, &createdUser)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	if err := database.Connect(); err != nil {
		panic(err.Error())
	}

	app := routes.New()
	type args struct {
		data        user.UpdateUser
		login       user.LoginUser
		userID      uint
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid update", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				ReachedTimeLimit: true,
				RemainingTime:    0,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:      updatedUser.ID,
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"User not found", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				ReachedTimeLimit: true,
				RemainingTime:    0,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:      updatedUser.ID + 1,
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				ReachedTimeLimit: true,
				RemainingTime:    0,
			},
			userID: updatedUser.ID,
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode: http.StatusBadRequest,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginBody, _ := json.Marshal(tt.args.login)
			reqLogin, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(loginBody))
			reqLogin.Header.Set("Content-Type", "application/json")

			var login = new(user.ResponseAuth)
			resLogin, _ := app.Test(reqLogin, -1)
			defer resLogin.Body.Close()
			resBodyLogin, _ := ioutil.ReadAll(resLogin.Body)
			json.Unmarshal(resBodyLogin, &login)

			reqBody, _ := json.Marshal(tt.args.data)
			endpoint := fmt.Sprintf("/api/v1/users/%d", tt.args.userID)
			req, _ := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
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
				Email:    "dinopuguh@ymail.com",
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

func TestDelete(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	r := routes.New()

	type args struct {
		userID     uint
		statusCode int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid delete user", args{
			userID:     createdUser.ID,
			statusCode: http.StatusOK,
		}},
		{"User not found", args{
			userID:     createdUser.ID + 2,
			statusCode: http.StatusNotFound,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := fmt.Sprintf("/api/v1/users/%d", tt.args.userID)
			req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)

			res, _ := r.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}
