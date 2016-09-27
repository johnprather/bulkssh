package bulkssh

// Result object resulting from Runner.OutCh
type Result struct {
	Request   *Request
	Outputs   map[*Command]string
	CmdErrors map[*Command]error
	Error     error
}
