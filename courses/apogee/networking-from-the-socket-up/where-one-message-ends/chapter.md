---
slug: where-one-message-ends
title: Where One Message Ends
position: 1
---
Frame messages on a stream: delimited lines with bufio.Reader, exact byte counts with io.ReadFull, length prefixes, one buffered reader per connection and one write per frame.
