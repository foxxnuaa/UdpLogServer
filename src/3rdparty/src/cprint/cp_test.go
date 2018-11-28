package cprint

import (
	_ "errors"
	"testing"
)

func TestCP(t *testing.T) {

	// test one model
	//P(DEBUG, "Test DEBUG color P() model", "\n")
	//P(WARING, "Remote latest version %v %v latest version %v.\n", "0.10.28", "=", "0.1.0.26")
	//P(DEBUG, "\n")

	//// test tow model
	//P(DEBUG, "Test custom color P() model", "\n")
	//cp := CP{Red, false, None, false, "="}
	//P(NOTICE, "Remote latest version %v %v latest version %v, don't need to upgrade.\n", "0.10.28", cp, "0.1.0.26")
	//P(DEBUG, "\n")

	//// test three model
	//P(DEBUG, "Test Error() model", "\n")
	//err := errors.New("Variable is not defined.")
	//Error(ERROR, "'gnvm updte latest' an error has occurred. \nError: ", err)
	//P(DEBUG, "\n")

	P(DEBUG, "Test DEBUG color P() model", "\n")
	P(WARING, "Test WARING color P() model", "\n")
	P(ERROR, "Test ERROR color P() model", "\n")
	P(NOTICE, "Test NOTICE color P() model", "\n")
}
