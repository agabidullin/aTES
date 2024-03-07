package events

const TaskCreated = "task_created"

type TaskCreatedPayload struct {
	PublicId    uint
	Title       string
	Description string
	AssigneeId  uint
}

const TaskAssigned = "task_assigned"

type TaskAssignedPayload struct {
	PublicId   uint
	AssigneeId uint
}

const TaskCompleted = "task_completed"

type TaskCompletedPayload struct {
	PublicId uint
}
