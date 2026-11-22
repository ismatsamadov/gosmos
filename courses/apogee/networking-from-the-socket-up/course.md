---
slug: networking-from-the-socket-up
title: Networking from the Socket Up
kind: Course
position: 5
short: Talk to Ystmark's new dish controllers byte by byte, then build a web server from a bare socket.
---
Ystmark's new dish controllers arrive with a manual that is wrong on page one, and Apogee drives them from Port Ellery across a line that drops in every storm. With Sunniva at the station and Nkem reviewing, you write the code that talks to them: connections and listeners, messages framed on a stream, big-endian fields and checksums, deadlines on every read and write, many connections at once, and reconnecting without sending a command twice. Then you build the station server from nothing but a socket: requests, bodies, chunked uploads and answers that Go's own client reads. Standard library only, and every connection stays on your Mac.
