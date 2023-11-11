package login

type User struct {
	Data    []string
	RetData []string
	Mark    bool
}

// NewUser creates a User.
func NewUser(data []string, _type int) (*User, error) {

	if _type == TypeLogin {
		if len(data) != 3 {
			return nil, ErrorInternalImpossibleInputData
		}
		user := &User{Data: data}
		user.RetData = make([]string, 3)

		return user, nil
	}

	if _type == TypeRegister {
		if len(data) != 4 {
			return nil, ErrorInternalImpossibleInputData
		}
		user := &User{Data: data}
		user.RetData = make([]string, 3)

		return user, nil
	}

	return nil, ErrorInternalImpossibleInputData
}
