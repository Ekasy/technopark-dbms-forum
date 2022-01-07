package error

import "fmt"

type CustomError struct {
	Code    int
	Message string
}

func (ce CustomError) Error() string {
	return fmt.Sprintf("error: code: %d message: %s", ce.Code, ce.Message)
}

func GenerateError(code int, message string) *CustomError {
	return &CustomError{
		Code: code, Message: message,
	}
}

var (
	InternalDbError CustomError = CustomError{
		Code:    500,
		Message: "internal database error",
	}

	RollbackError CustomError = CustomError{
		Code:    500,
		Message: "rollback databse error",
	}

	CommitError CustomError = CustomError{
		Code:    500,
		Message: "commit database error",
	}

	NicknameAlreadyExist CustomError = CustomError{
		Code:    500,
		Message: "user with this nickname already exist",
	}

	EmailAlreadyExist CustomError = CustomError{
		Code:    500,
		Message: "user with this email already exist",
	}

	NoRows CustomError = CustomError{
		Code:    404,
		Message: "no rows in result set",
	}

	ConflictData CustomError = CustomError{
		Code:    409,
		Message: "",
	}

	ForumAlreadyExist CustomError = CustomError{
		Code:    500,
		Message: "forum with this slug already exist",
	}

	UserNotExist CustomError = CustomError{
		Code:    500,
		Message: "user not exist",
	}

	AuthorNotExist CustomError = CustomError{
		Code:    500,
		Message: "forum with this author not exist",
	}

	ForumNotExist CustomError = CustomError{
		Code:    500,
		Message: "forum with this slug not exist",
	}

	ThreadAlreadyExist CustomError = CustomError{
		Code:    500,
		Message: "thread with this slug not exist",
	}

	ThreadNotExists CustomError = CustomError{
		Code:    500,
		Message: "thread not exist",
	}

	ParentNotExist CustomError = CustomError{
		Code:    500,
		Message: "parent's post not exist",
	}

	PostNotExist CustomError = CustomError{
		Code:    500,
		Message: "post not exist",
	}
)
