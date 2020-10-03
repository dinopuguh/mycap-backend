package group_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/mycap-backend/services/user"

	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/routes"
	"github.com/dinopuguh/mycap-backend/services/group"
	"github.com/stretchr/testify/assert"
)

var (
	createdGroup *group.Group
)

func TestNew(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	app := routes.New()

	type args struct {
		data        group.CreateGroup
		login       user.LoginUser
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Group type not specified", args{
			data: group.CreateGroup{
				Type: "Chat Room",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: group.CreateGroup{
				Type: "Group",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode: http.StatusBadRequest,
		}},
		{"Valid create group", args{
			data: group.CreateGroup{
				Type: "Group",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Group already exists", args{
			data: group.CreateGroup{
				Type: "Group",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
		}},
		{"Reached time limit", args{
			data: group.CreateGroup{
				Type: "Group",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusBadRequest,
			contentType: "application/json",
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				json.Unmarshal(resBody, &createdGroup)
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
	t.Run("Get all groups", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/groups", nil)

		res, _ := router.Test(req, -1)
		defer res.Body.Close()
		resBody, _ := ioutil.ReadAll(res.Body)

		assert.Equalf(t, http.StatusOK, res.StatusCode, string(resBody))
	})
}

func TestJoin(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	app := routes.New()

	type args struct {
		data        group.JoinGroup
		login       user.LoginUser
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Valid join group", args{
			data: group.JoinGroup{
				AdminUsername: "dinopuguh",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Group not found.", args{
			data: group.JoinGroup{
				AdminUsername: "anyusername99",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: group.JoinGroup{
				AdminUsername: "dinopuguh",
			},
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/join-groups", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))
		})
	}
}

func TestLeave(t *testing.T) {
	err := database.Connect()
	if err != nil {
		panic("Can't connect database.")
	}
	app := routes.New()

	type args struct {
		data        group.LeaveGroup
		login       user.LoginUser
		statusCode  int
		contentType string
	}
	tests := []struct {
		name string
		args args
	}{
		{"Group not found.", args{
			data: group.LeaveGroup{
				AdminUsername: "anyusername99",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusNotFound,
			contentType: "application/json",
		}},
		{"Body parser invalid", args{
			data: group.LeaveGroup{
				AdminUsername: "dinopuguh",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode: http.StatusBadRequest,
		}},
		{"Valid leave group", args{
			data: group.LeaveGroup{
				AdminUsername: "dinopuguh",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@email.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
		}},
		{"Valid end group", args{
			data: group.LeaveGroup{
				AdminUsername: "dinopuguh",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			statusCode:  http.StatusOK,
			contentType: "application/json",
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/leave-groups", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)

			assert.Equalf(t, tt.args.statusCode, res.StatusCode, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				endpoint := fmt.Sprintf("/api/v1/users/%d", login.User.ID)
				reqDeleteUser, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
				app.Test(reqDeleteUser, -1)
			}

		})
	}
}
