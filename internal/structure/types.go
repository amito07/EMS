package structure

type Student struct {
	Id    int
	Name  string `validate:"required"`
	Age   int	`validate:"required"`
	Email string `validate:"required,email"`
}

type StudentSignupBody struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
}

type StudentLoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type StudentLoginResponse struct {
	FirstName string `json:"first_name"`
	Token string `json:"token"`
}
