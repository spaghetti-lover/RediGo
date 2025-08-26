package core

import (
	"errors"
	"syscall"
)

func cmdPING(args []string) []byte {
	var res []byte

	// edge case
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), true)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}

	return res
}

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte

	switch cmd.Cmd { 
	case "PING":
		res = cmdPING(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}

	_, err := syscall.Write(connFd, res)
	return err
}
