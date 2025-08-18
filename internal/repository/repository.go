package repository

type Repository interface {
	ToDoRepository
	UsersRepository
}

type ToDoRepository interface {
	GetTaskUser()
	CreateTask()
	UpdateTask()
	CompleteTask()
	DeleteTask()
}

type UsersRepository interface {
	CreateUser()
	GetUsers()
	DeleteUser()
}
