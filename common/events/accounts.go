package events

const AccountCreated = "account_created"

type AccountCreatedPayload struct {
	PublicId uint
	Name     string
	Role     string
}

const RoleChanged = "role_changed"

type RoleChangedPayload struct {
	PublicId uint
	Role     string
}
