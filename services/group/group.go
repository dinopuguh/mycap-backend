package group

import (
	"net/http"

	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/dinopuguh/mycap-backend/database"
	"github.com/dinopuguh/mycap-backend/response"
	"github.com/dinopuguh/mycap-backend/services/user"
	"github.com/gofiber/fiber/v2"
)

// Group is a model for group chat
type Group struct {
	gorm.Model
	AdminID       uint        `json:"admin_id"`
	AdminUsername string      `json:"admin_username"`
	Admin         user.User   `json:"admin"`
	Type          string      `json:"type"`
	Participants  []user.User `json:"participants" gorm:"many2many:group_participants;"`
}

const (
	// GroupType is an enum for group chat
	GroupType = "Group"
	// ConferenceType is an enum for conference
	ConferenceType = "Conference"
)

// GetAll is a function to get all groups from database
// @Summary Get all groups
// @Description Get all groups
// @Tags groups
// @Accept json
// @Produce json
// @Success 200 {object} response.HTTP
// @Failure 200 {object} response.HTTP
// @Router /v1/groups [get]
func GetAll(c *fiber.Ctx) error {
	db := database.DBConn

	var groups []Group
	if res := db.Preload("Admin").Preload("Participants").Find(&groups); res.Error != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data:    groups,
		Status:  http.StatusOK,
		Message: "Success get all groups.",
	})
}

// New function creates a group/conference for communicate with deaf people
// @Summary Create a group chat or conference
// @Description Create a group chat or conference
// @Tags groups
// @Accept json
// @Produce json
// @Param group body CreateGroup true "Create group"
// @Success 200 {object} response.HTTP
// @Failure 200 {object} response.HTTP
// @Failure 401 {object} string
// @Security ApiKeyAuth
// @Router /v1/groups [post]
func New(c *fiber.Ctx) error {
	db := database.DBConn

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var admin = new(user.User)
	db.Where("email = ?", email).First(&admin)

	var existingGroup = new(Group)
	if res := db.Where("admin_id = ?", admin.ID).First(&existingGroup); res.RowsAffected != 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: "You are already has a group chat or conference.",
		})
	}

	if admin.ReachedTimeLimit {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: "This user already reached time limit this month.",
		})
	}

	var group = new(Group)

	participants := make([]user.User, 0)
	participants = append(participants, *admin)

	createGroup := new(CreateGroup)
	if err := c.BodyParser(&createGroup); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	group.Admin = *admin
	group.AdminUsername = admin.Username
	group.Participants = participants

	if createGroup.Type != GroupType && createGroup.Type != ConferenceType {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: "Group type not specified.",
		})
	}

	group.Type = createGroup.Type

	if res := db.Create(group); res.Error != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusServiceUnavailable,
			Message: res.Error.Error(),
		})
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data:    group,
		Status:  http.StatusOK,
		Message: "Success create a new group.",
	})
}

// Join function assign an user to specific group
// @Summary Joining group chat or conference
// @Description Joining group chat or conference
// @Tags groups
// @Accept json
// @Produce json
// @Param group body JoinGroup true "Join group"
// @Success 200 {object} response.HTTP
// @Failure 200 {object} response.HTTP
// @Failure 401 {object} string
// @Security ApiKeyAuth
// @Router /v1/join-groups [post]
func Join(c *fiber.Ctx) error {
	db := database.DBConn

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var joiningUser = new(user.User)
	db.Where("email = ?", email).First(&joiningUser)

	joinGroup := new(JoinGroup)
	if err := c.BodyParser(&joinGroup); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var group = new(Group)
	if res := db.Where("admin_username = ?", joinGroup.AdminUsername).First(&group); res.RowsAffected == 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusNotFound,
			Message: "Group not found.",
		})
	}

	group.Participants = append(group.Participants, *joiningUser)

	db.Save(&group)

	db.Preload("Admin").Preload("Participants").First(&group, group.ID)

	return c.JSON(response.HTTP{
		Success: true,
		Data:    group,
		Status:  http.StatusOK,
		Message: "Success joining a group.",
	})
}

// Leave function remove an user from specific group
// @Summary Leaving group chat or conference
// @Description Leaving group chat or conference
// @Tags groups
// @Accept json
// @Produce json
// @Param group body LeaveGroup true "Leave group"
// @Success 200 {object} response.HTTP
// @Failure 200 {object} response.HTTP
// @Failure 401 {object} string
// @Security ApiKeyAuth
// @Router /v1/leave-groups [post]
func Leave(c *fiber.Ctx) error {
	db := database.DBConn

	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	var leavingUser = new(user.User)
	db.Where("email = ?", email).First(&leavingUser)

	leaveGroup := new(LeaveGroup)
	if err := c.BodyParser(&leaveGroup); err != nil {
		return c.JSON(response.HTTP{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	var group = new(Group)
	if res := db.Preload("Admin").Where("admin_username = ?", leaveGroup.AdminUsername).First(&group); res.RowsAffected == 0 {
		return c.JSON(response.HTTP{
			Status:  http.StatusNotFound,
			Message: "Group not found.",
		})
	}

	if group.AdminID == leavingUser.ID {
		leavingUser.RemainingTime = leaveGroup.RemainingTime
		if leaveGroup.RemainingTime == 0 {
			leavingUser.ReachedTimeLimit = true
		}
		db.Save(&leavingUser)

		db.Model(&group).Association("Participants").Clear()
		db.Delete(&group)
	} else {
		db.Model(&group).Association("Participants").Delete(leavingUser)
		db.Preload("Admin").Preload("Participants").First(&group, group.ID)
	}

	return c.JSON(response.HTTP{
		Success: true,
		Data:    group,
		Status:  http.StatusOK,
		Message: "Success leaving group.",
	})
}
