package users

import "fmt"

// User represents a user entity.
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// IUserService defines the user service interface.
type IUserService interface {
	GetUser(id int) (*User, error)
	ListUsers() ([]*User, error)
	CreateUser(user *User) error
}

// UserService is the default implementation of IUserService.
type UserService struct {
	users map[int]*User
}

// NewUserService creates a new UserService.
func NewUserService() IUserService {
	return &UserService{
		users: map[int]*User{
			1: {ID: 1, Name: "Alice", Email: "alice@example.com"},
			2: {ID: 2, Name: "Bob", Email: "bob@example.com"},
		},
	}
}

func (s *UserService) GetUser(id int) (*User, error) {
	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	return user, nil
}

func (s *UserService) ListUsers() ([]*User, error) {
	users := make([]*User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users, nil
}

func (s *UserService) CreateUser(user *User) error {
	if _, exists := s.users[user.ID]; exists {
		return fmt.Errorf("user already exists: %d", user.ID)
	}
	s.users[user.ID] = user
	return nil
}

