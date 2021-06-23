package main

import (
	"context"
	"fmt"
	"strings"
)

type OpHandler struct {
	sess *Session
}

func NewOpHandler(s *Session) *OpHandler {
	return &OpHandler{s}
}

func (h *OpHandler) Handle(ctx context.Context, cmd string, param []byte) ([]byte, error) {
	if strings.HasPrefix(cmd, "/") {
		// normal mode
		var (
			out []byte
			err error
		)
		switch cmd {
		case "/remember":
			fallthrough
		case "/r":
			out, err = h.handleRemember(ctx, param)
		case "/ask":
			fallthrough
		case "/a":
			out, err = h.handleAsk(ctx, param)
		default:
			out = []byte(h.helpMsg())
		}
		if err != nil {
			return nil, err
		}
		return []byte(out), nil
	} else if strings.HasPrefix(cmd, "%") {
		// admin mode
	}
	return []byte(h.helpMsg()), nil
}

func (h *OpHandler) handleRemember(ctx context.Context, param []byte) ([]byte, error) {
	id, err := PutNewTip(ctx, string(param))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("Cmd: /remember Param: %s Status: SAVED RECORD_ID: %s Session: %s", string(param), string(id), h.sess.ID)), nil
}

func (h *OpHandler) handleAsk(ctx context.Context, param []byte) ([]byte, error) {
	out, err := SearchTip(ctx, string(param))
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("Cmd: /ask Param: %s Output: %s Session: %s", string(param), out, h.sess.ID)), nil
}

func (h *OpHandler) helpMsg() string {
	return fmt.Sprintf("Usage:\n" +
		"/remember or /r <text> : save new knowledge.\n" +
		"/ask or /a <keyword> : search knowledge by keyword.")
}
