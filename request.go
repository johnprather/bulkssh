package bulkssh

// Request object for sending to Runner.InCh
type Request struct {
	Hostname string
	Port     int
	Commands []*Command
	Result   *Result
}

// AddCommand - use to append commands to the request
func (r *Request) AddCommand(command string) {
	newcmd := &Command{
		command,
	}
	r.Commands = append(r.Commands, newcmd)
}
