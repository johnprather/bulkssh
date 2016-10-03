package bulkssh

// Request object for sending to Runner.InCh
type Request struct {
	Hostname       string
	Port           int
	User           *string
	Password       *string
	Agent          *bool
	Commands       []*Command
	Error          error
	ConnectTimeout int
	CommandTimeout int
}

// AddCommand - use to append commands to the request
func (r *Request) AddCommand(command string) {
	newcmd := &Command{}
	newcmd.Command = command
	r.Commands = append(r.Commands, newcmd)
}

// NewRequest - generates a new request object, with user,
// host, and port populated, which needs only have commands
// added before it can be submitted
func NewRequest(user string, host string, port int) *Request {
	req := &Request{}
	req.User = &user
	req.Hostname = host
	req.Port = port
	req.Commands = make([]*Command, 0)
	req.ConnectTimeout = 10
	req.CommandTimeout = 0
	return req
}
