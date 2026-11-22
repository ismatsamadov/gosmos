---
slug: a-client-that-gives-up
title: A Client That Gives Up
type: code_tests
position: 0
difficulty: 5
topics: http, time
created: 2026-11-07T06:18:31Z
modified: 2026-11-07T06:18:31Z
---
Last Tuesday the go/no-go poll, the checks before every pass, asked the weather service, Tomás's program that reports each station's wind, about Cape Wren and waited four minutes. Furrow-B came into view, the moment known as AOS, before its checklist was done, and Ingrid wants the poll's client to give up in time.

`http.Get` uses `http.DefaultClient`, which waits as long as the server takes. A client of your own can have a `Timeout`: any call through it that runs longer, reading the body included, fails with an error:

```go
c := &http.Client{Timeout: 3 * time.Second}
resp, err := c.Get(url)
if err != nil {
	return err // after 3 seconds at most
}
defer resp.Body.Close()
```

## Your task

`newWeatherClient` gives the weather client `&http.Client{}`, which never gives up. Give it a client of its own whose `Timeout` is `timeout`. `Reading` ignores `ctx` until the next lesson.

> [!WARNING]
> Don't set `http.DefaultClient.Timeout` instead: every part of a program shares that client.
