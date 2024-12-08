package shared_utils

import "golang.org/x/crypto/bcrypt"

func Filter[T any](arr []T, predicate func(T) bool) []T {
	var result []T
	for _, item := range arr {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](input []T, mapper func(T) R) []R {
	var result []R
	for _, item := range input {
		result = append(result, mapper(item))
	}
	return result
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckHashPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
