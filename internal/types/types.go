package types

type Student struct {
	Id       string
	Name     string `validate:"required"`
	Email    string `validate:"required"`
	Age      int    `validate:"required"`
	Password string `validate:"required"`
}
