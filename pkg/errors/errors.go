package errors

import "fmt"

type DelayError struct {
	Format string
	Message string
}

func ( de DelayError) Error() string {
	return fmt.Sprintf(de.Format,de.Message)
}

func ( de DelayError) Errorf(format , message string) string {
	return fmt.Sprintf(format,message)
}