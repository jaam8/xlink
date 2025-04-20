package utils

func GetRoleByIsStaffIsAdmin(isStaff bool, isAdmin bool) string {
	if isStaff && isAdmin {
		return "superuser"
	}
	if isAdmin {
		return "admin"
	}
	if isStaff {
		return "staff"
	}
	return "user"
}
