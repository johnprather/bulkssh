package bulkssh

// Runner the main ssh runner object of the bulkssh package
type Runner struct {
	InCh  chan *Request
	OutCh chan *Request
}

// NewRunner instantiates a Runner object
func NewRunner(routines int) *Runner {
	inch := make(chan *Request)
	outch := make(chan *Request)
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
	connection, err := sshInit(req)
	if err != nil {
		req.Error = err
	} else {
		defer sshDisconnect(connection)
		for _, cmd := range req.Commands {
			output, err := sshRun(connection, cmd, req.CommandTimeout)
			cmd.Error = err
			if output != nil {
				cmd.Output = *output
			}
		}
	}

	r.OutCh <- req
}
