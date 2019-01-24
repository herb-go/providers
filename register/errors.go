package register

var Debug = false

type DuplicationError struct {
	Type RegisterType
	Key  string
}

func (e *DuplicationError) Error() string {
	return "register error: \"" + e.Key + "\" has registered to type \"" + string(e.Type) + "\""
}

func IsDuplicationError(err error) bool {
	_, ok := err.(*DuplicationError)
	return ok
}

type NotRegsiteredError struct {
	Type RegisterType
	Key  string
}

func (e *NotRegsiteredError) Error() string {
	return "register error: \"" + e.Key + "\"  (type \"" + string(e.Type) + "\") is not registered "
}

func IsNotRegsiteredError(err error) bool {
	_, ok := err.(*DuplicationError)
	return ok
}
