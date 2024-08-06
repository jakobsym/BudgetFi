package budgetfi

import (
	"context"
	"errors"

	"github.com/jakobsym/BudgetFi/api/pkg/model"
)

var ErrNotFound = errors.New("not found")

type budgetfiRepo interface {
	Post(ctx context.Context, user *model.User) error
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

func (c *Controller) Post(ctx context.Context, user *model.User) error {
	return c.repo.Post(ctx, user)
}

//TODO: Get
//TODO: Put
//TODO: Delete
