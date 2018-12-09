package domain

type Document struct {
	ID   string
	Name string
	Size int64
}

const PathFolder = "C:\\Users\\UTI01\\Downloads\\Programs"

type Result struct {
	Code int
	Message string
	Json string
}

type RequestBase struct {
	Action string
}

type RequestDeleteDocument struct {
	Base RequestBase
	ID string
}

type RequestGetDocument struct {
	Base RequestBase
	ID string
}


type RequestDownloadDocument struct {
	Base RequestBase
	ID string
}


type RequestListDocument struct {
	Base RequestBase
}


type RequestCreateDocument struct {
	Base RequestBase
	Name string
	Content []byte
}