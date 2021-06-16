# Jarvis

A Websocket-based Knowledge Base Robot Demo


Usage:

```
/remember (or /r) <Text> : Add new knowledge text
/ask (or /a) <Keyword> : Search by keyword
```

Example:

Server:
```
$ go build && ./javis
```

Client:

```
$ websocat ws://localhost:8080/ws
/r hello world
Cmd: /remember Param: hello world Status: SAVED RECORD_ID: 7 Session: 3a2740ce-139e-4e1d-ac14-7b4ff83029d5
/a hello
Cmd: /ask Param: hello Output: KB-7 hello world  Session: 3a2740ce-139e-4e1d-ac14-7b4ff83029d5
/r bye world
Cmd: /remember Param: bye world Status: SAVED RECORD_ID: 8 Session: 3a2740ce-139e-4e1d-ac14-7b4ff83029d5
/a world
Cmd: /ask Param: world Output: KB-7 hello world KB-8 bye world  Session: 3a2740ce-139e-4e1d-ac14-7b4ff83029d5

```
