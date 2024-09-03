package cli

type Args struct {
	TemplateId string
	FileList   string
	Threads    int
	Timeout    int
	MatchOnly  bool
	Verbose    bool
	Show       bool
	All        bool
}

var currentArgs Args
