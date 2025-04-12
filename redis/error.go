package redis

import "fmt"

func ErrHSet(err error) error {
	return fmt.Errorf("failed to HSet due to error: %v", err)
}

func ErrHGet(err error) error {
	return fmt.Errorf("failed to HGet due to error: %v", err)
}

func ErrHGetAll(err error) error {
	return fmt.Errorf("failed to HGetAll due to error: %v", err)
}

func ErrHDel(err error) error {
	return fmt.Errorf("failed to HDel due to error: %v", err)
}

func ErrHExists(err error) error {
	return fmt.Errorf("failed to HExists due to error: %v", err)
}
