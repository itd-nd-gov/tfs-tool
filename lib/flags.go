package lib

type flagsT struct {
	UserID         string
	Password       string
	DestinationDir string
	Verbose        bool
	Color          bool
}

var Flags flagsT
