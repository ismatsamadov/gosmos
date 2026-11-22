Keep the newest frames in one map from satellite ID to `TelemetryFrame`, behind a `sync.Mutex`: many passes take in telemetry at the same moment. A frame takes the place of the one kept only when its `Seq` is higher.

Write the report once every contact has ended, only when `ctx` hasn't: walk the network's `Satellites` in order, write `formatLine` of the kept frame or the never-heard line for each, and save the file whole, as the plan's files are saved.

If a line holds an older frame than the tests want, the frame kept was the last to arrive, not the highest-numbered: a satellite over Ystmark and Tirrin Rock at once comes in out of order when Ystmark relays late. If Thistle-2's line reads `nan`, a line reached the map without passing `parseLine`.
