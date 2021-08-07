package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"github.com/Egor-Ruban/PIP/graph/generated"
	"github.com/Egor-Ruban/PIP/graph/model"
)

func (r *queryResolver) Object(ctx context.Context, id *string) (*model.Object, error) {
	object := model.Object{}

	const q = `SELECT * FROM objects WHERE object_id = $1`

	if err := r.DB.Get(&object, q, id); err != nil {
		return nil, err
	}

	return &object, nil
}

func (r *queryResolver) Subject(ctx context.Context, id *string) (*model.Subject, error) {
	subject := model.Subject{}

	const q = `SELECT * FROM subjects WHERE subject_id = $1`

	if err := r.DB.Get(&subject, q, id); err != nil {
		return nil, err
	}

	return &subject, nil
}

func (r *queryResolver) Environment(ctx context.Context) (*model.Environment, error) {
	time := time.Now()
	env := model.Environment{Date: &time}

	return &env, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
