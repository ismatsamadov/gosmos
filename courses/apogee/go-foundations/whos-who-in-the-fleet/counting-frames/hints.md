Every ID in `ids` stands for one frame. What has to happen to that satellite's count each time its ID comes up?

Loop over `ids` with `range`, and add one to the entry in `counts` under the current ID with `++`. No check is needed for a satellite you haven't seen yet.
