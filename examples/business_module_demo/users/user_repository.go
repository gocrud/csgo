package users

// IUserRepository defines the user repository interface.
type IUserRepository interface {
	FindByID(id int) (*User, error)
	FindAll() ([]*User, error)
	Save(user *User) error
}

// UserRepository is the default implementation of IUserRepository.
type UserRepository struct {
	// In real application, this would be a database connection
	store map[int]*User
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository() IUserRepository {
	return &UserRepository{
		store: make(map[int]*User),
	}
}

func (r *UserRepository) FindByID(id int) (*User, error) {
	user, ok := r.store[id]
	if !ok {
		return nil, nil
	}
	return user, nil
}

func (r *UserRepository) FindAll() ([]*User, error) {
	users := make([]*User, 0, len(r.store))
	for _, u := range r.store {
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepository) Save(user *User) error {
	r.store[user.ID] = user
	return nil
}

