package messages

const RegisterAccountKey = "register-account"

type RegisterAccount struct {
	PublicId uint
	Name     string
	Role     string
}

const ChangeRoleKey = "change-role"

type ChangeRole struct {
	PublicId uint
	Role     string
}
