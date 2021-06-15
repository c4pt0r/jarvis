package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"time"

	"github.com/ngaut/log"
)

var CONTEXT_VAL_SESSION = "session"

type Session struct {
	ctx        context.Context
	ID         string
	rdr        io.ReadCloser
	wr         io.WriteCloser
	createAt   time.Time
	lastAction time.Time

	done chan *Session

	sm *StateMachine
}

func NewSession(ctx context.Context,
	ID string,
	rd io.ReadCloser, w io.WriteCloser) *Session {
	ret := &Session{
		ID:       ID,
		ctx:      ctx,
		rdr:      rd,
		wr:       w,
		createAt: time.Now(),
		sm:       NewStateMachine(),
		done:     make(chan *Session),
	}
	return ret
}

func (s *Session) Close() {
	log.Info("session is closing...")
	err := s.rdr.Close()
	if err != nil {
		log.Error(err)
	}
	s.wr.Close()
	c := context.WithValue(s.ctx, CONTEXT_VAL_SESSION, s)
	s.sm.OnClose(c)

	s.done <- s
	close(s.done)
}

func (s *Session) DoneChan() chan *Session {
	return s.done
}

func (s *Session) Loop() {
	c := context.WithValue(s.ctx, CONTEXT_VAL_SESSION, s)
	s.sm.OnStart(c)
	rdr := bufio.NewReader(s.rdr)
	var err error
	for {
		var line, ret []byte
		var exit bool
		line, _, err = rdr.ReadLine()
		if err != nil {
			break
		}
		// format: command params
		parts := bytes.SplitN(line, []byte(" "), 2)
		cmd := parts[0]
		var param []byte
		if len(parts) > 1 {
			param = parts[1]
		}
		ret, exit, err = s.sm.OnOp(cmd, param)
		if err != nil {
			break
		}
		s.lastAction = time.Now()
		_, err = s.wr.Write(ret)
		if err != nil || exit {
			s.wr.Write([]byte("Bye!"))
			break
		}
	}
	if err != nil {
		log.Error(err)
		s.wr.Write([]byte(err.Error()))
	}
	s.Close()
}
