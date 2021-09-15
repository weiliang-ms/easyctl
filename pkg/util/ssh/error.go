package ssh

import "fmt"

type Waitmsg struct {
	status int
	signal string
	msg    string
	lang   string
}

type ExitError struct {
	Waitmsg
}

// ExitStatus returns the exit status of the remote command.
func (w Waitmsg) ExitStatus() int {
	return w.status
}

func (e *ExitError) Error() string {
	return e.Waitmsg.String()
}

func (w Waitmsg) String() string {
	str := fmt.Sprintf("Process exited with status %v", w.status)
	if w.signal != "" {
		str += fmt.Sprintf(" from signal %v", w.signal)
	}
	if w.msg != "" {
		str += fmt.Sprintf(". Reason was: %v", w.msg)
	}
	return str
}
