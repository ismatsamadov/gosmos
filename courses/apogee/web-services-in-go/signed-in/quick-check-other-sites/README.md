---
slug: quick-check-other-sites
title: Quick Check: Other Sites
type: quiz
position: 6
difficulty: 8
topics: http, security
created: 2026-11-13T03:49:39Z
modified: 2026-11-13T12:18:58Z
---
Before the page goes to the beta, Nkem has three questions about requests that other sites' pages start. This is how the service now meets them:

```text
address     caller shows   guard
/account/   a cookie       noCrossOrigin
the API     a token        cors(one origin)
```

## Your task

Answer three questions: why only the account page needs a guard against other sites' forms, what CORS does about a form posted from another site, and what `Vary: Origin` is for.
