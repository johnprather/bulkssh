package bulkssh

// Command a ssh command
type Command struct {
	Command string
	Error   error
	Output  string
}
