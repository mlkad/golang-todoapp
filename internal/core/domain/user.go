package domain

type User struct {
	ID 			int
	Version int
	
	Fullname string
	PhoneNumber *string
} 

//GET /users/{id}