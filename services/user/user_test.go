package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/response"
	"github.com/dinopuguh/mycap-backend/routes"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/stretchr/testify/assert"
)

var (
	createdUser *user.User
	updatedUser *user.User
)

func TestNew(t *testing.T) {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	app := routes.New()

	type args struct {
		data          user.RegisterUser
		expectDBError bool
		statusCode    int
		contentType   string
		willUpdate    bool
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
				TypeID:   1,
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
				TypeID:   1,
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
				TypeID:   1,
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
			willUpdate:  true,
		}},
		{"User's type not found", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dinopuguh16@email.com",
				Username: "dino1603",
				Password: "s3cr3tp45sw0rd",
				TypeID:   99,
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
			willUpdate:  true,
		}},
		{"Email already exist.", args{
			data: user.RegisterUser{
				Name:     "Dino",
				Email:    "dino@email.com",
				Username: "dinopuguh16",
				Password: "12345678",
				TypeID:   1,
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
				TypeID:   1,
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
				TypeID:   1,
			},
			statusCode: http.StatusBadRequest,
		}},
		{"DB connection closed", args{
			data: user.RegisterUser{
				Name:     "Dino Puguh",
				Email:    "dinopuguh@mycap9.com",
				Username: "dinopuguh9",
				Password: "s3cr3tp45sw0rd",
				TypeID:   1,
			},
			expectDBError: true,
			statusCode:    http.StatusServiceUnavailable,
			contentType:   "application/json",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.expectDBError {
				db, _ := database.DBConn.DB()
				db.Close()
			}

			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)

			resHTTP := new(response.HTTP)
			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				auth := new(user.ResponseAuth)
				authJSON, _ := json.Marshal(resHTTP.Data)
				json.Unmarshal(authJSON, &auth)

				if tt.args.willUpdate {
					updatedUser = &auth.User
				} else {
					createdUser = &auth.User
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
		data          user.UpdateUser
		login         user.LoginUser
		expectDBError bool
		userID        uint
		statusCode    int
		contentType   string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid update", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				RemainingTime:    1800,
				ReachedTimeLimit: false,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:      updatedUser.ID,
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Valid update 2", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				RemainingTime:    0,
				ReachedTimeLimit: true,
				TypeID:           1,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:      updatedUser.ID,
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"User type not found", args{
			data: user.UpdateUser{
				Name:   "Dino Yang Baru",
				TypeID: 99,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:      updatedUser.ID,
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"User not found", args{
			data: user.UpdateUser{
				Name: "Dino Yang Baru",
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
				Name: "Dino Yang Baru",
			},
			userID: updatedUser.ID,
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode: http.StatusBadRequest,
		}},
		{"DB connection closed", args{
			data: user.UpdateUser{
				Name:             "Dino Yang Baru",
				RemainingTime:    0,
				ReachedTimeLimit: true,
				TypeID:           1,
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			expectDBError: true,
			userID:        updatedUser.ID,
			statusCode:    http.StatusServiceUnavailable,
			contentType:   "application/json",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginBody, _ := json.Marshal(tt.args.login)
			reqLogin, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(loginBody))
			reqLogin.Header.Set("Content-Type", "application/json")

			resHTTP := new(response.HTTP)
			login := new(user.ResponseAuth)
			resLogin, _ := app.Test(reqLogin, -1)
			defer resLogin.Body.Close()
			resBodyLogin, _ := ioutil.ReadAll(resLogin.Body)
			json.Unmarshal(resBodyLogin, &resHTTP)
			loginJSON, _ := json.Marshal(resHTTP.Data)
			json.Unmarshal(loginJSON, &login)

			if tt.args.expectDBError {
				db, _ := database.DBConn.DB()
				db.Close()
			}

			reqBody, _ := json.Marshal(tt.args.data)
			endpoint := fmt.Sprintf("/api/v1/users/%d", tt.args.userID)
			req, _ := http.NewRequest(http.MethodPut, endpoint, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))
		})
	}
}

func TestGetAll(t *testing.T) {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	app := routes.New()

	type args struct {
		expectDBError bool
		statusCode    int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid get all", args{
			expectDBError: false,
			statusCode:    http.StatusOK,
		}},
		{"DB connection closed", args{
			expectDBError: true,
			statusCode:    http.StatusServiceUnavailable,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.expectDBError {
				db, _ := database.DBConn.DB()
				db.Close()
			}

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)

			resHTTP := new(response.HTTP)
			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))
		})
	}
}

func TestLogin(t *testing.T) {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	app := routes.New()

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

			resHTTP := new(response.HTTP)
			res, _ := app.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))
		})
	}
}

func TestDelete(t *testing.T) {
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	app := routes.New()

	type args struct {
		login         user.LoginUser
		expectDBError bool
		userID        uint
		statusCode    int
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid delete user", args{
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:     createdUser.ID,
			statusCode: http.StatusOK,
		}},
		{"User not found", args{
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			userID:     createdUser.ID + 2,
			statusCode: http.StatusNotFound,
		}},
		{"DB connection closed", args{
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			expectDBError: true,
			userID:        createdUser.ID + 1,
			statusCode:    http.StatusServiceUnavailable,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loginBody, _ := json.Marshal(tt.args.login)
			reqLogin, _ := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(loginBody))
			reqLogin.Header.Set("Content-Type", "application/json")

			resHTTP := new(response.HTTP)
			login := new(user.ResponseAuth)
			resLogin, _ := app.Test(reqLogin, -1)
			defer resLogin.Body.Close()
			resBodyLogin, _ := ioutil.ReadAll(resLogin.Body)
			json.Unmarshal(resBodyLogin, &resHTTP)
			loginJSON, _ := json.Marshal(resHTTP.Data)
			json.Unmarshal(loginJSON, &login)

			if tt.args.expectDBError {
				db, _ := database.DBConn.DB()
				db.Close()
			}

			endpoint := fmt.Sprintf("/api/v1/users/%d", tt.args.userID)
			req, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))
		})
	}
}
