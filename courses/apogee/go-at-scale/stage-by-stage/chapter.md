---
slug: stage-by-stage
title: Stage by Stage
position: 0
---
Build pipelines of goroutines joined by channels: a stage owns and closes the channel it returns, stages compose, every send and receive also watches the context, a sink drains the end, a fan-in joins many feeds, a stage's error reaches the whole run as a cancel cause, and a stage runs on several goroutines without losing its order.
