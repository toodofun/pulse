package checker

import (
	"fmt"
	"pulse/internal/checker/http"
	"pulse/internal/model"
)

type Checker interface {
	Check(fields string) *model.Record
	Validate(fields string) error
}

func GetChecker(t model.CheckerType) (Checker, error) {
	switch t {
	case http.CheckerTypeHTTP:
		return &http.Checker{}, nil
	default:
		return nil, fmt.Errorf("unknown checker: %s", t)
	}
}
