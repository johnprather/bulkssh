package bulkssh

// Runner the main ssh runner object of the bulkssh package
type Runner struct {
	InCh  chan *Request
	OutCh chan *Result
}

// NewRunner instantiates a Runner object
func NewRunner(routines int) *Runner {
	inch := make(chan *Request)
	outch := make(chan *Result)
	runner := &Runner{
		inch,
		outch,
	}

	for i := 0; i < routines; i++ {
		go runner.requestListener()
	}

	return runner
}

func (r *Runner) requestListener() {
	//fmt.Printf("dataFetcher-%d launched\n", id)
	for {
		select {
		case req := <-r.InCh:
			//fmt.Printf("dataFetcher-%d starting request for %s\n", id, req.host.hostname)
			r.handleRequest(req)
		}
	}
}

func (r *Runner) handleRequest(req *Request) {
	result := &Result{}
	result.Outputs = make(map[*Command]string)
	result.CmdErrors = make(map[*Command]error)
	result.Request = req
	req.Result = result

	connection, err := sshInit(req.Hostname, req.Port)
	if err != nil {
		result.Error = err
	} else {
		for _, cmd := range req.Commands {
			output, err := sshRun(connection, cmd)
			result.Outputs[cmd] = *output
			result.CmdErrors[cmd] = err
		}
	}
	r.OutCh <- result
}
