package dto_error

import "fmt"

func ErrAuthorizeFor(resource string) error {
	return fmt.Errorf("you are not authorized for %s", resource)
}