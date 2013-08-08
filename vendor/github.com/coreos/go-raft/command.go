package raft

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

var commandTypes map[string]Command

func init() {
	commandTypes = map[string]Command{}
}

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A command represents an action to be taken on the replicated state machine.
type Command interface {
	CommandName() string
	Apply(server *Server) (interface{}, error)
}

type CommandEncoder interface {
	Encode(w io.Writer) error
	Decode(r io.Reader) error
}

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

//--------------------------------------
// Instantiation
//--------------------------------------

// Creates a new instance of a command by name.
func newCommand(name string, data []byte) (Command, error) {
	// Find the registered command.
	command := commandTypes[name]
	if command == nil {
		return nil, fmt.Errorf("raft.Command: Unregistered command type: %s", name)
	}

	// Make a copy of the command.
	v := reflect.New(reflect.Indirect(reflect.ValueOf(command)).Type()).Interface()
	copy, ok := v.(Command)
	if !ok {
		panic(fmt.Sprintf("raft: Unable to copy command: %s (%v)", command.CommandName(), reflect.ValueOf(v).Kind().String()))
	}

	// If data for the command was passed in the decode it.
	if data != nil {
		if encoder, ok := copy.(CommandEncoder); ok {
			if err := encoder.Decode(bytes.NewReader(data)); err != nil {
				return nil, err
			}
		} else {
			json.NewDecoder(bytes.NewReader(data)).Decode(copy)
		}
	}

	return copy, nil
}

//--------------------------------------
// Registration
//--------------------------------------

// Registers a command by storing a reference to an instance of it.
func RegisterCommand(command Command) {
	if command == nil {
		panic(fmt.Sprintf("raft: Cannot register nil"))
	} else if commandTypes[command.CommandName()] != nil {
		panic(fmt.Sprintf("raft: Duplicate registration: %s", command.CommandName()))
		return
	}
	commandTypes[command.CommandName()] = command
}
