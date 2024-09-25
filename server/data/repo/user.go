package repo

import (
	"github.com/ginx-contribs/ginx-server/server/data/ent"
	"github.com/ginx-contribs/ginx-server/server/data/ent/user"
	"golang.org/x/net/context"
)

func NewUserRepo(client *ent.Client) *UserRepo {
	return &UserRepo{Ent: client}
}

type UserRepo struct {
	Ent *ent.Client
}

// FindByUID returns a User matching the given uid
func (u *UserRepo) FindByUID(ctx context.Context, uid string) (*ent.User, error) {
	return u.Ent.User.
		Query().
		Where(
			user.UIDEQ(uid),
		).Only(ctx)
}

// FindByNameOrMail returns a User matching the given name or email
func (u *UserRepo) FindByNameOrMail(ctx context.Context, name string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.Or(
				user.UsernameEQ(name),
				user.EmailEQ(name),
			),
		).Only(ctx)
}

// FindByName returns a user matching the given name
func (u *UserRepo) FindByName(ctx context.Context, name string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.UsernameEQ(name),
		).Only(ctx)
}

// FindByEmail returns a User matching the given email
func (u *UserRepo) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return u.Ent.User.Query().
		Where(
			user.EmailEQ(email),
		).Only(ctx)
}

// CreateNewUser creates a new user with the minimum information
func (u *UserRepo) CreateNewUser(ctx context.Context, username string, email string, password string) (*ent.User, error) {
	return u.Ent.User.Create().
		SetUsername(username).
		SetEmail(email).
		SetPassword(password).
		Save(ctx)
}

func (u *UserRepo) RemoveByUID(ctx context.Context, uid string) (int, error) {
	return u.Ent.User.Delete().
		Where(user.UIDEQ(uid)).
		Exec(ctx)
}

// UpdateOnePassword updates the user password with specified email
func (u *UserRepo) UpdateOnePassword(ctx context.Context, id int, password string) (*ent.User, error) {
	return u.Ent.User.UpdateOneID(id).
		SetPassword(password).
		Save(ctx)
}

// ListByPage list users by page
func (u *UserRepo) ListByPage(ctx context.Context, page, size int, pattern string) ([]*ent.User, error) {
	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	}

	query := u.Ent.User.Query()

	// pattern
	if pattern != "" {
		query = query.
			Where(
				user.Or(
					user.UsernameContains(pattern),
					user.EmailContains(pattern),
				),
			)
	}

	// pagination
	users, err := query.Offset((page - 1) * size).Limit(size).All(ctx)
	if err != nil {
		return []*ent.User{}, nil
	}
	return users, nil
}
