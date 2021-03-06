package utils

import "fmt"

// MultiError will convert a slice of errors into a single error
func MultiError(errs []error) error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		text := "Multiple Errors:\n"
		for _, err := range errs {
			text += err.Error() + "\n"
		}

		return fmt.Errorf(text)
	}
}
