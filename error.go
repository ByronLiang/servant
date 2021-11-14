package servant

import "fmt"

type netError struct {
	cause error
	kind  string
}

func (n *netError) Error() string {
	return fmt.Sprintf("[%s]: %s", n.kind, n.cause.Error())
}

func IsNetError(err error) bool {
	if _, ok := err.(*netError); ok {
		return true
	}
	return false
}
