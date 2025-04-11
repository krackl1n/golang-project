package apperr

import "github.com/pkg/errors"

var ErrorNotFound = errors.New(
	"not found",
)
