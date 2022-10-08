package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

const dbTimeout = time.Second * 3

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: User{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	User User
}

// User is the structure which holds one user from the database.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users order by last_name`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var users []*User
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			if err.Error() == sql.ErrNoRows.Error() {
				return nil, errors.New("no user found")
			}
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (u *User) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users where email=$1`
	row := db.QueryRowContext(ctx, query, email)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) GetOne(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `select id, email, first_name, last_name, password, user_active, created_at, updated_at
	from users where email=$1`
	row := db.QueryRowContext(ctx, query, id)
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) Update() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `update users set
	email = $1,
	first_name = $2,
	last_name = $3,
	user_active = $4,
	updated_at = $5
	where id = $6`
	_, err := db.ExecContext(ctx, query, u.Email, u.FirstName, u.LastName, u.Active, time.Now(), u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) Delete() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	query := `delete from users where id = $1`
	_, err := db.ExecContext(ctx, query, u.ID)
	if err != nil {
		return nil
	}
	return nil
}

func (u *User) Insert(user User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return 0, err
	}
	var newID int
	query := `insert into user into users (email, first_name, last_name, password, user_active, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6, $7) returning id`
	err = db.QueryRowContext(ctx, query, u.Email,
		u.FirstName,
		u.LastName,
		hashedPassword,
		u.Active,
		time.Now(),
		time.Now()).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (u *User) ResetPassword(password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	query := `update users set password = $1 where id = $2`
	_, err = db.ExecContext(ctx, query, hashedPassword, u.ID)
	if err != nil {
		return err
	}
	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
	if err != nil {
		switch err.Error() {
		case bcrypt.ErrMismatchedHashAndPassword.Error():
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}
