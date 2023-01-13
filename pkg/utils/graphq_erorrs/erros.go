package graphq_erorrs

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func InternalServerError(ctx context.Context) error {
	return &gqlerror.Error{
		Message: "internal server error",
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": 500,
		},
	}
}

func InternalServerErrorWithDesc(ctx context.Context, err error) error {
	return &gqlerror.Error{
		Message: fmt.Sprintf("internal server error: %v", err),
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": 500,
		},
	}
}

func CustomError(ctx context.Context, errCode int, message string) error {
	return &gqlerror.Error{
		Message: message,
		Path:    graphql.GetPath(ctx),
		Extensions: map[string]interface{}{
			"code": errCode,
		},
	}
}
