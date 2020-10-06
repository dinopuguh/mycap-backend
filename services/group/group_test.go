package group_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/dinopuguh/mycap-backend/response"
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
	if err := database.Connect(); err != nil {
		panic("Can't connect database.")
	}

	app := routes.New()

	type args struct {
		data          group.CreateGroup
		login         user.LoginUser
		expectDBError bool
		statusCode    int
		contentType   string
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
		{"DB connection closed", args{
			data: group.CreateGroup{
				Type: "Group",
			},
			login: user.LoginUser{
				Email:    "dinopuguh@mycap.com",
				Password: "s3cr3tp45sw0rd",
			},
			expectDBError: true,
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
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/groups", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				groupJSON, _ := json.Marshal(resHTTP.Data)
				json.Unmarshal(groupJSON, &createdGroup)
			}
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

			req, _ := http.NewRequest(http.MethodGet, "/api/v1/groups", nil)

			resHTTP := new(response.HTTP)
			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))
		})
	}
}

func TestJoin(t *testing.T) {
	if err := database.Connect(); err != nil {
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

			resHTTP := new(response.HTTP)
			login := new(user.ResponseAuth)
			resLogin, _ := app.Test(reqLogin, -1)
			defer resLogin.Body.Close()
			resBodyLogin, _ := ioutil.ReadAll(resLogin.Body)
			json.Unmarshal(resBodyLogin, &resHTTP)
			loginJSON, _ := json.Marshal(resHTTP.Data)
			json.Unmarshal(loginJSON, &login)

			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/join-groups", bytes.NewBuffer(reqBody))
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

func TestLeave(t *testing.T) {
	if err := database.Connect(); err != nil {
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

			resHTTP := new(response.HTTP)
			login := new(user.ResponseAuth)
			resLogin, _ := app.Test(reqLogin, -1)
			defer resLogin.Body.Close()
			resBodyLogin, _ := ioutil.ReadAll(resLogin.Body)
			json.Unmarshal(resBodyLogin, &resHTTP)
			loginJSON, _ := json.Marshal(resHTTP.Data)
			json.Unmarshal(loginJSON, &login)

			reqBody, _ := json.Marshal(tt.args.data)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/leave-groups", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", tt.args.contentType)
			req.Header.Set("Authorization", "Bearer "+login.AccessToken)

			res, _ := app.Test(req, -1)
			defer res.Body.Close()
			resBody, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &resHTTP)

			assert.Equalf(t, tt.args.statusCode, resHTTP.Status, string(resBody))

			if tt.args.statusCode == http.StatusOK {
				endpoint := fmt.Sprintf("/api/v1/users/%d", login.User.ID)
				reqDeleteUser, _ := http.NewRequest(http.MethodDelete, endpoint, nil)
				reqDeleteUser.Header.Set("Authorization", "Bearer "+login.AccessToken)
				app.Test(reqDeleteUser, -1)
			}

		})
	}
}
