package accounts

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type UpdateForm struct {
	Name  string `schema:"name"`
	Email string `schema:"email"`
}
