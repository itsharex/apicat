package api

import (
	"math"
	"net/http"
	"strings"

	"github.com/apicat/apicat/common/auth"
	"github.com/apicat/apicat/common/translator"
	"github.com/apicat/apicat/models"
	"github.com/gin-gonic/gin"
)

type GetMembersData struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"page_size" binding:"omitempty,gte=1,lte=100"`
}

type AddMemberData struct {
	Email    string `json:"email" binding:"required,email,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=255"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

type UserIDData struct {
	UserID int `uri:"user-id" binding:"required,gte=1"`
}

type SetMemberData struct {
	Email     string `json:"email" binding:"omitempty,email,lte=255"`
	Password  string `json:"password" binding:"omitempty,gte=6,lte=255"`
	Role      string `json:"role" binding:"omitempty,oneof=admin user"`
	IsEnabled int    `json:"is_enabled" binding:"required,oneof=0 1"`
}

func GetMembers(ctx *gin.Context) {
	var data GetMembersData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindQuery(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data.Page <= 0 {
		data.Page = 1
	}
	if data.PageSize <= 0 {
		data.PageSize = 15
	}

	user, _ := models.NewUsers()
	totalUsers, err := user.Count()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.QueryFailed"}),
		})
		return
	}

	users, err := user.List(data.Page, data.PageSize)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.QueryFailed"}),
		})
		return
	}

	userList := []any{}
	if len(users) > 0 {
		for _, v := range users {
			email := v.Email
			parts := strings.Split(email, "@")

			userList = append(userList, map[string]any{
				"id":         v.ID,
				"username":   v.Username,
				"email":      parts[0][0:1] + "***" + parts[0][len(parts[0])-1:] + "@" + parts[len(parts)-1],
				"role":       v.Role,
				"is_enabled": v.IsEnabled,
				"created_at": v.CreatedAt.Format("2006-01-02 15:04:05"),
				"updated_at": v.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"current_page": data.Page,
		"total_page":   int(math.Ceil(float64(totalUsers) / float64(data.PageSize))),
		"total_member": totalUsers,
		"members":      userList,
	})
}

func AddMember(ctx *gin.Context) {
	var data AddMemberData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, _ := models.NewUsers()
	user.Email = data.Email
	if err := user.GetByEmail(data.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	hashedPassword, err := auth.HashPassword(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.CreateFailed"}),
		})
		return
	}

	user.Username = strings.Split(data.Email, "@")[0]
	user.Role = data.Role
	user.Password = hashedPassword

	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.CreateFailed"}),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"role":       user.Role,
		"is_enabled": user.IsEnabled,
		"created_at": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updated_at": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	})
}

func SetMember(ctx *gin.Context) {
	var userIDData UserIDData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&userIDData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	var data SetMemberData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindJSON(&data)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.NewUsers(uint(userIDData.UserID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
		})
		return
	}

	if err := user.GetByEmail(data.Email); err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "User.MailboxAlreadyExists"}),
		})
		return
	}

	if data.Email != "" {
		user.Email = data.Email
	}
	if data.Role != "" {
		user.Role = data.Role
	}
	if data.IsEnabled != 0 {
		user.IsEnabled = data.IsEnabled
	}
	if data.Password != "" {
		hashedPassword, err := auth.HashPassword(data.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
			})
			return
		}
		user.Password = hashedPassword
	}

	if err := user.Save(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.UpdateFailed"}),
		})
		return
	}

	ctx.Status(http.StatusCreated)
}

func DeleteMember(ctx *gin.Context) {
	var userIDData UserIDData
	if err := translator.ValiadteTransErr(ctx, ctx.ShouldBindUri(&userIDData)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	user, err := models.NewUsers(uint(userIDData.UserID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	if err := user.Delete(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": translator.Trasnlate(ctx, &translator.TT{ID: "Member.DeleteFailed"}),
		})
		return
	}

	ctx.Status(http.StatusNoContent)
}
