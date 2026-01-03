package shared

import (
	"errors"
	"net/mail"
	"time"
	"unicode"
)

var (
	ErrKeyword            = "error"
	ErrWrongInputs        = "you have entered wrong input"
	ErrBadRequest         = "user bad request"
	ErrNotFoundUser       = "user is not found"
	ErrWrongCredentials   = "invalid username or password"
	ErrWrongPassword      = "invalid password"
	ErrListingUsers       = "cannot find any users"
	ErrAddingUser         = "new user cannot be added"
	ErrUpdatingUser       = "user cannot be updated"
	ErrDeletingUser       = "user cannot be deleted"
	ErrAddingIncome       = "income cannot be added"
	ErrUpdatingIncome     = "income cannot be updated"
	ErrDeleteIncome       = "income cannot be deleted"
	ErrListingIncomes     = "incomes cannot be listed"
	ErrAddingExpense      = "expense cannot be added"
	ErrUpdatingExpense    = "expense cannot be updated"
	ErrDeletingExpense    = "expense cannot be deleted"
	ErrListingExpenses    = "expenses cannot be listed"
	ErrJWTTokenGenFailed  = "JWT token generation failed"
	ErrWrongRole          = "Wrong role selected"
	ErrInternalError      = "Internal Error"
	ErrMissedRoleProperty = "role property is not found in token"
	ErrOnlyAdmin          = "only admin can do this"
	ErrAddingCategory     = "category cannot be added"
	ErrUpdatingCategory   = "category cannot be updated"
	ErrDeleteCategory     = "category cannot be deleted"
	ErrListingCategories  = "categories cannot be listed"
	ErrPayingExpense      = "upcoming expense cannot be moved"

	OkKeyword          = "ok"
	OkSignupUser       = "user signed up successfully"
	OkAddingUser       = "user added successfully"
	OkUpdatingUser     = "user updated successfully"
	OkDeletingUser     = "user deleted successfully"
	OkAddingIncome     = "income added successfully"
	OkUpdatingIncome   = "income updated successfully"
	OkDeletingIncome   = "income deleted successfully"
	OkAddingExpense    = "expense added successfully"
	OkUpdatingExpense  = "expense updated successfully"
	OkDeletingExpense  = "expense deleted successfully"
	OkPassChanged      = "password changed successfully"
	TokenKey           = "token"
	UserLogoutMsg      = "user logged out successfully"
	OkAddingCategory   = "category added successfully"
	OkUpdatingCategory = "category updated successfully"
	OkDeletingCategory = "category deleted successfully"
	OkPayingExpense    = "upcoming expenses moved to expenses successfully"
)

var (
	ErrTooShort     = errors.New("password is too short")
	ErrNoUppercase  = errors.New("password must contain at least one uppercase letter")
	ErrNoLowercase  = errors.New("password must contain at least one lowercase letter")
	ErrNoDigit      = errors.New("password must contain at least one digit")
	ErrNoSpecial    = errors.New("password must contain at least one special character")
	ErrInvalidEmail = errors.New("email is not valid")
)

var Roles = []string{"admin", "user"}

// UnixTimeToRFC339 used to convert time from unixtime to RFC3339 - like 2006-01-02T15:04:05Z07:00
func UnixTimeToRFC339(t uint64) string {
	return time.Unix(int64(t), 0).Format(time.RFC3339)
}

const PasswordMinLength = 8

func ValidatePassword(pw string) error {
	if len(pw) < PasswordMinLength {
		return ErrTooShort
	}

	var upper, lower, digit, special bool

	for _, ch := range pw {
		switch {
		case unicode.IsUpper(ch):
			upper = true
		case unicode.IsLower(ch):
			lower = true
		case unicode.IsDigit(ch):
			digit = true
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			special = true
		}
	}

	if !upper {
		return ErrNoUppercase
	}
	if !lower {
		return ErrNoLowercase
	}
	if !digit {
		return ErrNoDigit
	}
	if !special {
		return ErrNoSpecial
	}

	return nil
}

func IsValidEmail(email string) error {
	_, err := mail.ParseAddress(email)

	if err != nil {
		return ErrInvalidEmail
	}

	return nil
}
