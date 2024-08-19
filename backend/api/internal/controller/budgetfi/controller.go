package budgetfi

import (
	"context"
	"errors"

	"github.com/jakobsym/BudgetFi/api/pkg/model"
)

var ErrNotFound = errors.New("not found")

type budgetfiRepo interface {
	//Post(ctx context.Context, user *model.User) error
	CreateUser(ctx context.Context, user *model.User) error
	PrevUserCheck(ctx context.Context, user *model.User) (string, error)
	//TODO: Get
	//TODO: Put
	//TODO: Delete
}

type Controller struct {
	repo budgetfiRepo
}

func New(repo budgetfiRepo) *Controller {
	return &Controller{repo: repo}
}

/*
// Deprecated

	func (c *Controller) Post(ctx context.Context, user *model.User) error {
		return c.repo.Post(ctx, user)
	}
*/
func (c *Controller) PrevUserCheck(ctx context.Context, user *model.User) (string, error) {
	return c.repo.PrevUserCheck(ctx, user)
}

func (c *Controller) CreateUser(ctx context.Context, user *model.User) error {
	return c.repo.CreateUser(ctx, user)
}

//TODO: Get
//TODO: Put
//TODO: Delete
