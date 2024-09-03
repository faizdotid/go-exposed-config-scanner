package cli

type Args struct {
	Show       bool
	All        bool
	Threads    int
	TemplateId string
	FileList   string
	MatchOnly  bool
	Verbose    bool
	Timeout    int
}

var currentArgs Args
