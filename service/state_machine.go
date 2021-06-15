package main

import (
	"context"
	"strings"

	"github.com/ngaut/log"
)

var (
	StateMachineStoreSessionKey string = "session"
)

type StateMachine struct {
	store map[string]interface{}
	// set onStart
	sess      *Session
	opHandler *OpHandler
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		store: make(map[string]interface{}),
	}
}

func (sm *StateMachine) sessionInfo() string {
	if sm.sess != nil {
		return sm.sess.ID
	}
	return ""
}

func (sm *StateMachine) OnStart(ctx context.Context) error {
	sess := ctx.Value(CONTEXT_VAL_SESSION).(*Session)
	// set session
	sm.sess = sess
	// init handler set session
	sm.opHandler = NewOpHandler(sess)
	log.Info("State machine onStart...", "session:", sm.sessionInfo())
	return nil
}

func (sm *StateMachine) OnClose(ctx context.Context) error {
	log.Info("State machine onClose...", "session:", sm.sessionInfo())
	return nil
}

func (sm *StateMachine) OnOp(op []byte, param []byte) (result []byte, exit bool, err error) {
	log.Info("State machine Op:", op, "session:", sm.sessionInfo())
	cmd := strings.ToLower(string(op))
	if cmd == "exit" {
		return []byte("Bye"), true, nil
	}
	var ret []byte
	ret, err = sm.opHandler.Handle(cmd, param)
	if err != nil {
		log.Error(err)
		ret = []byte(err.Error())
	}
	return ret, false, nil
}
