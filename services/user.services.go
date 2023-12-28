package services

import (
	"github.com/emarifer/go-echo-templ-htmx/db"
	"golang.org/x/crypto/bcrypt"
)

func NewUserServices(u User, uStore db.Store) *UserServices {

	return &UserServices{
		User:      u,
		UserStore: uStore,
	}
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserServices struct {
	User      User
	UserStore db.Store
}

func (us *UserServices) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 8)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users(email, password, username) VALUES($1, $2, $3)`

	_, err = us.UserStore.Db.Exec(
		stmt,
		u.Email,
		string(hashedPassword),
		u.Username,
	)

	return err
}

func (us *UserServices) CheckEmail(email string) (User, error) {

	query := `SELECT id, email, password, username FROM users
		WHERE email = ?`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Email = email
	err = stmt.QueryRow(
		us.User.Email,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}

/* func (us *UserServices) GetUserById(id int) (User, error) {

	query := `SELECT id, email, password, username FROM users
		WHERE id = ?`

	stmt, err := us.UserStore.Db.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.ID = id
	err = stmt.QueryRow(
		us.User.ID,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
} */
