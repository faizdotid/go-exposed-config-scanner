package cli

type Args struct {
	TemplateId string
	List       string
	Threads    int
	Timeout    int
	MatchOnly  bool
	Verbose    bool
	Show       bool
	All        bool
}

var currentArgs Args
