package tests

import (
	"testing"
	"xlink/user_service/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestGetRoleByIsStaffIsAdmin(t *testing.T) {
	t.Run("superuser case", func(t *testing.T) {
		isStaff := true
		isAdmin := true
		resp := utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

		assert.Equal(t, "superuser", resp)
	})

	t.Run("admin case", func(t *testing.T) {
		isStaff := false
		isAdmin := true
		resp := utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

		assert.Equal(t, "admin", resp)
	})

	t.Run("staff case", func(t *testing.T) {
		isStaff := true
		isAdmin := false
		resp := utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

		assert.Equal(t, "staff", resp)
	})

	t.Run("user case", func(t *testing.T) {
		isStaff := false
		isAdmin := false
		resp := utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

		assert.Equal(t, "user", resp)
	})

}
