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

type AvatarForm struct {
	UserID   string `schema:"userID"`
	Filename string `schema:"fileName"`
}

type ProfileForm struct {
	Name   string `schema:"name"`
	Email  string `schema:"email"`
	Avatar AvatarForm
}
