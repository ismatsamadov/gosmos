/// Invented inputs for lesson-style jq queries, set in the Apogee world of docs/apogee: the
/// logs Apogee's pass scheduler writes, the reports its check programs print, and the
/// responses of its booking API.
enum Fixture {
    // MARK: Logs, read as JSON lines

    /// The scheduler's log over one morning at Port Ellery: start-up with build info, a pass
    /// planned, calls handled (one with the operator redacted), an uplink refused with its
    /// stack, a wind warning, a lookup of a station that does not exist, and a failed booking.
    static let schedulerLog = #"""
        {"time":"2027-03-01T06:00:02.114Z","level":"INFO","msg":"scheduler started","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","port":7340}
        {"time":"2027-03-01T06:00:05.902Z","level":"INFO","msg":"pass planned","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","satellite":"pipit-1","pass_id":"p-0301-pel-01","aos":"2027-03-01T06:41:10Z","los":"2027-03-01T06:52:48Z","max_elevation":61.4}
        {"time":"2027-03-01T06:41:12.330Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-first-light","method":"GET","path":"/","client":"10.20.0.x","duration_us":412,"request_bytes":0,"response_bytes":96,"status":200}
        {"time":"2027-03-01T06:41:40.870Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-0f3a51","method":"POST","path":"/v1/bookings","client":"10.20.0.x","duration_us":2210,"request_bytes":118,"response_bytes":204,"status":201,"operator":"[REDACTED]"}
        {"time":"2027-03-01T06:44:58.001Z","level":"ERROR","msg":"uplink refused","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-77b2c4","satellite":"pipit-2","operator":"furrowlight","error":{"message":"signature does not match the command key","stack":["uplink.(*Sender).Verify\n\t/srv/apogee/uplink/verify.go:73","uplink.(*Sender).Send\n\t/srv/apogee/uplink/send.go:41","net/http.HandlerFunc.ServeHTTP\n\t/usr/local/go/src/net/http/server.go:2322"]}}
        {"time":"2027-03-01T06:44:58.004Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-77b2c4","method":"POST","path":"/v1/commands","client":"10.20.0.x","duration_us":904,"request_bytes":512,"response_bytes":61,"status":403,"error":{"message":"command refused"}}
        {"time":"2027-03-01T06:52:50.010Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-9c40aa","method":"GET","path":"/v1/passes","client":"10.20.0.x","duration_us":15422,"request_bytes":0,"response_bytes":3880,"status":200,"operator":"[REDACTED]"}
        {"time":"2027-03-01T06:58:13.551Z","level":"WARN","msg":"tracking in wind","station":"CWR","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","wind_kmh":57}
        {"time":"2027-03-01T07:03:40.118Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-1d0e93","method":"GET","path":"/v1/stations/XXX","client":"10.20.0.x","duration_us":311,"request_bytes":0,"response_bytes":28,"status":404,"error":{"message":"no such station"}}
        {"time":"2027-03-01T07:05:22.907Z","level":"INFO","msg":"handled call","station":"PEL","hostname":"pel-sched-1","build":"7c1e9a2","build_time":"2027-02-28T21:40:00Z","trace_id":"tr-5e18f0","method":"POST","path":"/v1/bookings","client":"10.20.0.x","duration_us":5001377,"request_bytes":121,"response_bytes":36,"status":500,"error":{"message":"pass table is locked","stack":["store.(*Passes).Reserve\n\t/srv/apogee/store/passes.go:118","api.bookPass\n\t/srv/apogee/api/bookings.go:56"]}}

        """#

    /// The log of an earlier lesson, before errors were objects: the stack is part of the text.
    static let textStackLog = #"""
        {"time":"2027-03-02T09:12:40.551Z","level":"ERROR","msg":"uplink refused","satellite":"pipit-2","error":"signature does not match the command key:\n    uplink.(*Sender).Verify\n        /srv/apogee/uplink/verify.go:73\n    uplink.(*Sender).Send\n        /srv/apogee/uplink/send.go:41"}
        {"time":"2027-03-02T09:12:40.553Z","level":"INFO","msg":"handled call","method":"POST","path":"/v1/commands","status":403}
        """#

    /// After asking for a downlink file the archive does not have.
    static let missingFileLog = #"""
        {"time":"2027-03-03T14:20:44.913Z","level":"ERROR","msg":"downlink file missing","pass_id":"p-0303-sqt-04","error":"open archive/p-0303-sqt-04/frames.bin: no such file or directory"}
        {"time":"2027-03-03T14:20:44.914Z","level":"INFO","msg":"handled call","method":"GET","path":"/v1/downlinks/p-0303-sqt-04","status":404}
        """#

    /// After listing downlinks with one unreadable entry in the archive.
    static let listFailedLog = #"""
        {"time":"2027-03-03T14:31:07.210Z","level":"ERROR","msg":"listing downlinks failed","error":{"message":"read archive/CORRUPT: is a directory","path":"archive/CORRUPT"}}
        {"time":"2027-03-03T14:31:07.211Z","level":"INFO","msg":"handled call","method":"GET","path":"/v1/downlinks","status":500}
        """#

    /// The same with two unreadable entries, logged as joined errors.
    static let joinedErrorsLog = #"""
        {"time":"2027-03-03T14:36:52.004Z","level":"ERROR","msg":"listing downlinks failed","errors":{"error_1":{"message":"read archive/TORN-B: is a directory","path":"archive/TORN-B"},"error_2":{"message":"read archive/TORN-A: is a directory","path":"archive/TORN-A"}}}
        {"time":"2027-03-03T14:36:52.005Z","level":"INFO","msg":"handled call","method":"GET","path":"/v1/downlinks","status":500}
        """#

    /// The audit log of the operator portal, with personal data redacted where it must be.
    static let auditLog = #"""
        {"time":"2027-03-04T11:02:13.481Z","level":"INFO","event":"sign_in","email":"[REDACTED]","operatorId":null,"success":false,"reason":"unknown_key"}
        {"time":"2027-03-04T11:02:19.006Z","level":"INFO","event":"sign_in","email":"[REDACTED]","operatorId":14,"success":true}
        {"time":"2027-03-04T11:03:40.770Z","level":"INFO","event":"booking_paid","email":"[REDACTED]","operatorId":14,"bookingId":3301,"totalCents":129000,"billingName":"[REDACTED]","billingAddress":"[REDACTED]","billingCity":"[REDACTED]","billingCountry":"[REDACTED]","billingPostcode":"[REDACTED]"}
        {"time":"2027-03-04T11:05:02.118Z","level":"INFO","event":"key_reset_request","email":"ops@furrowlight.example","success":true,"resetLink":"http://localhost:7340/key-reset/9b2e41c07d3f5a86e1c4b9d02f7a3e5c8b1d6f0a2c4e7b9d1f3a5c8e0b2d4f6a"}
        {"time":"2027-03-04T11:06:31.902Z","level":"INFO","event":"console_opened","operatorId":9,"email":"ops@furrowlight.example","expiresAt":"2027-04-03T11:06:31Z"}
        """#

    // MARK: Command output, read as one JSON value

    /// `go env -json GOVERSION`, as Go prints it.
    static func goEnv(_ version: String) -> String {
        "{\n\t\"GOVERSION\": \"\(version)\"\n}\n"
    }

    /// The report of a check program over a revoked command key.
    static let keyRevocationReport = #"{"revoked":2,"oldKeyAccepted":false,"graceKeyAccepted":false,"newKeyAccepted":true,"revokedAtMatchesClock":true,"secondRevocationNoChange":true}"#

    /// The report of a check program over a pass's timing.
    static let passTimingReport = #"{"aosBeforeLos":true,"durationSeconds":698,"maxElevationAboveMask":true,"withinBooking":true}"#

    /// The report of a check program that bursts calls at the API: the status of each call.
    static let burstReport: String = {
        func statuses(_ status: Int, _ count: Int, then last: Int) -> String {
            "[" + (Array(repeating: String(status), count: count) + [String(last)]).joined(separator: ",") + "]"
        }
        return """
            {"bookingStatuses":\(statuses(201, 4, then: 429)),"stationStatuses":\(statuses(200, 30, then: 429)),\
            "commandStatuses":\(statuses(202, 12, then: 429)),"retryAfterOnEveryRefusal":true,"refusalBodiesMatch":true}
            """
    }()

    // MARK: HTTP responses: the booking API

    static let operatorID = "5b0f8e2a-41c3-4d7e-9a6b-2c8d1f4e7a90"
    static let otherOperatorID = "c71d3a95-6e2b-4f08-b4c1-9d7e5a2f3b16"
    static let bookingID = "0e9a4c27-3b1d-4f6a-8e52-7c4b9d1a6f38"

    static let createdOperator = #"{"id":"\#(operatorID)","created_at":"2027-03-01T08:15:02.123456Z","updated_at":"2027-03-01T08:15:02.123456Z","email":"ops@furrowlight.example","priority":false}"#

    static let updatedOperator = #"{"id":"\#(operatorID)","created_at":"2027-03-01T08:15:02.123456Z","updated_at":"2027-03-01T08:19:47.554012Z","email":"control@furrowlight.example","priority":false}"#

    static let signIn = #"{"id":"\#(otherOperatorID)","created_at":"2027-03-01T08:16:40.902117Z","updated_at":"2027-03-01T08:16:40.902117Z","email":"duty@cirrine.example","priority":true,"token":"eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcG9nZWUtYm9va2luZyIsInN1YiI6ImM3MWQzYTk1IiwiZXhwIjoxODA0MDAwMDAwfQ.Zx4kQ9wT1mB7nR2vY8cL5pH3sJ6dF0gA","refresh_token":"4d8a1f6c3b9e2d7a0f5c8b1e4d7a2f9c6b3e0d8a5f2c9b6e3d0a7f4c1b8e5d2a"}"#

    static let booking = #"{"id":"\#(bookingID)","created_at":"2027-03-01T08:17:12.330981Z","updated_at":"2027-03-01T08:17:12.330981Z","pass_id":"p-0301-pel-01","satellite":"pipit-1","station":"PEL","operator_id":"\#(otherOperatorID)"}"#

    static func passes(_ satellites: [String]) -> String {
        "[" + satellites.enumerated().map { index, satellite in
            #"{"id":"p-0301-yst-0\#(index)","satellite":"\#(satellite)","station":"YST","aos":"2027-03-01T1\#(index):00:00Z","max_elevation":\#(20 + index * 10)}"#
        }.joined(separator: ",") + "]"
    }

    static let passesByTime = passes(["pipit-1", "furrow-a", "furrow-b", "cirrine-1"])
    static let passesByElevation = passes(["cirrine-1", "furrow-b", "furrow-a", "pipit-1"])

    static func station(_ code: String, mask: Double, dishes: Int, open: Bool, booked: Int, stats: String) -> String {
        #"{"code":"\#(code)","network":"apogee","bands":["S","X"],"mask_deg":\#(mask),"dishes":\#(dishes),"accepting":\#(open),"staffed":\#(code != "CWR"),"passes_booked":\#(booked),"passes_free":\#(booked / 2),"pass_stats":\#(stats),"timezone":"UTC"}"#
    }

    static let ystmark = station("YST", mask: 5, dishes: 2, open: true, booked: 61, stats: #"{"planned":61,"tracked":58,"delivered":57,"missed":3}"#)
    static let quietYstmark = station("YST", mask: 5, dishes: 2, open: true, booked: 0, stats: #"{"planned":61,"tracked":61,"delivered":61,"missed":0}"#)
    static let capeWren = station("CWR", mask: 5, dishes: 1, open: false, booked: 0, stats: #"{"planned":12,"tracked":12,"missed":0}"#)
    static let saltFlat = station("SQT", mask: 3, dishes: 2, open: true, booked: 7, stats: #"{"planned":7,"tracked":7,"missed":0}"#)
    static let atoll = station("LAN", mask: 8, dishes: 1, open: true, booked: 4, stats: #"{"planned":4,"missed":0}"#)
    static let stormedYstmark = station("YST", mask: 5, dishes: 2, open: true, booked: 0, stats: #"{"planned":44,"tracked":17,"delivered":15,"stowed":27,"missed":129}"#)

    static let bookedColumn = #"{"passes_booked":512}"#

    static let network = #"{"api_version":"2.4.0","rate_mode":"basic","product":"Apogee Ground","name":"apogee","region":"global","stations_online":4,"pass_stats":{"planned":1342,"planned_rate":{"per_hour":3.2},"refused":3,"refused_rate":{"per_hour":0.0},"delivered":1290},"totals":{"stations":5,"dishes":7,"operators":6,"satellites":8,"bookings":214}}"#

    static let polarRoute = #"{"name":"polar","network":"apogee","kind":"direct","stations":["YST","CWR"],"internal":false}"#
    static let equatorRoute = #"{"name":"equator","network":"apogee","kind":"fanout","stations":["LAN","SQT"],"internal":false}"#
    static let routeBindings = #"[{"route":"polar","network":"apogee","satellite":"pipit-1","key":"pipit","priority":2}]"#

    // MARK: HTTP responses: the operations dashboards

    static let boardSearch = #"[{"id":1,"uid":"apogee-next-passes","title":"The Board","url":"/d/apogee-next-passes/the-board","tags":["fishbowl"],"starred":false}]"#

    static let alertRules = #"""
        [{"id":1,"uid":"a71c2e90","folder":"ops","group":"Fishbowl","title":"Missed Pass Burst","for":"2m","data":[{"ref":"A","window":{"from":900,"to":0},"model":{"expr":"sum(increase(passes_missed_total[15m]))","ref":"A"}},{"ref":"C","window":{"from":0,"to":0},"model":{"type":"threshold","conditions":[{"evaluator":{"type":"gt","params":[2]},"reducer":{"type":"last"}}]}}]},{"id":2,"uid":"b82d3f01","folder":"ops","group":"Fishbowl","title":"Uplink Refusal Burst","for":"1m","data":[{"ref":"A","window":{"from":600,"to":0},"model":{"expr":"sum(rate(uplinks_refused_total[5m])) * 60","ref":"A"}},{"ref":"C","window":{"from":0,"to":0},"model":{"type":"threshold","conditions":[{"evaluator":{"type":"gt","params":[4]},"reducer":{"type":"last"}}]}}]}]
        """#

    static let firingRules = #"{"status":"success","data":{"groups":[{"name":"Fishbowl","rules":[{"state":"firing","name":"Missed Pass Burst","health":"ok"},{"state":"firing","name":"Uplink Refusal Burst","health":"ok"}]},{}]}}"#

    static let gaugeQuery = #"{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"dish_azimuth_degrees","station":"SQT","dish":"chiquita"},"value":[1803902410.512,"214.5"]}]}}"#

    // MARK: HTTP responses: downlink files

    static let fileID = "9d3e7b12-6a4f-4c80-b2e5-1f7a9c3d5e68"

    static func downlinks(preview: String, frames: String, telemetry: String) -> String {
        #"[{"id":"\#(fileID)","created_at":"2027-03-05T12:01:44.552Z","pass_id":"p-0305-yst-02","preview_url":"\#(preview)","file_url":"\#(frames)","title":"Furrow-A swath 0305-02","description":"Multispectral frames over the Ystmark fjords","operator_id":"\#(operatorID)"},{"id":"4b8e2c61-0d7f-4a93-8e1b-6c2f9a5d3b07","created_at":"2027-03-05T12:12:05.118Z","pass_id":"p-0305-yst-03","preview_url":"\#(preview)","file_url":"\#(telemetry)","title":"Pipit-1 beacon dump","description":"The echo of the first uplinked message","operator_id":"\#(operatorID)"}]"#
    }

    static let downlinksServedByAPI = downlinks(
        preview: "http://localhost:7340/v1/previews/\(fileID)",
        frames: "http://localhost:7340/archive/\(fileID).bin",
        telemetry: "http://localhost:7340/archive/4b8e2c61.bin"
    )
    static let downlinksWithInlinePreviews = downlinks(
        preview: "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
        frames: "http://localhost:7340/archive/\(fileID).bin",
        telemetry: "http://localhost:7340/archive/4b8e2c61.bin"
    )
    static let downlinksInBuckets = downlinks(
        preview: "http://localhost:7340/archive/Rk7mW2xP9bLs4nT8cQ6yHa3dF5gJ1vZ0.png",
        frames: "https://apogee-downlinks.storage.example/frames/aB3nQ8mK4xR6bW1cY9sL2dF7.bin",
        telemetry: "https://apogee-downlinks.storage.example/telemetry/hJ5gP0aZ8kU2eN6rX4wB1qM9.bin"
    )
    static let downlinksSigned = downlinks(
        preview: "http://localhost:7340/archive/Rk7mW2xP9bLs4nT8cQ6yHa3dF5gJ1vZ0.png",
        frames: "https://apogee-downlinks.storage.example/frames/aB3nQ8mK4xR6bW1cY9sL2dF7.bin?expires=1804000300&signature=6e1f3c8a0b5d2e7f9a4c1b6d8e0f2a5c7b9d1e3f5a7c9e1b3d5f7a9c1e3b5d7f",
        telemetry: "https://apogee-downlinks.storage.example/telemetry/hJ5gP0aZ8kU2eN6rX4wB1qM9.bin?expires=1804000300"
    )
    static let downlinksOnEdge = downlinks(
        preview: "https://edge.apogee.example/archive/Rk7mW2xP9bLs4nT8cQ6yHa3dF5gJ1vZ0.png",
        frames: "https://edge.apogee.example/frames/aB3nQ8mK4xR6bW1cY9sL2dF7.bin",
        telemetry: "https://edge.apogee.example/telemetry/hJ5gP0aZ8kU2eN6rX4wB1qM9.bin"
    )

    static let downlink = #"{"id":"\#(fileID)","created_at":"2027-03-05T12:01:44.552Z","pass_id":"p-0305-yst-02","preview_url":"http://localhost:7340/v1/previews/\#(fileID)","file_url":null,"title":"Furrow-A swath 0305-02","description":"Multispectral frames over the Ystmark fjords","operator_id":"\#(operatorID)"}"#

    // MARK: HTTP responses: the operator portal

    static let health = #"{"ok":true,"app":"operator-portal","version":"2.4.0"}"#

    /// A booking before the lesson that hides internal fields: Go's field names, no tags.
    static let bookingWithInternalFields = #"{"booking":{"ID":1,"OperatorEmail":"ops@furrowlight.example","BillingName":"Furrowlight Ltd","BillingAddress":"4 Quay Lane","InternalNotes":"pays late, book anyway","TotalCents":129000,"Status":"paid","CreatedAt":"2027-03-04T10:21:09Z"},"lines":[{"PassID":"p-0304-yst-01","Station":"YST","Minutes":11,"PriceCents":64500},{"PassID":"p-0304-sqt-02","Station":"SQT","Minutes":9,"PriceCents":64500}]}"#

    static let paidBooking = #"{"booking":{"id":1,"status":"paid","total_cents":129000,"created_at":"2027-03-04T10:21:09Z"},"lines":[{"pass_id":"p-0304-yst-01","station":"YST","minutes":11,"price_cents":64500},{"pass_id":"p-0304-sqt-02","station":"SQT","minutes":9,"price_cents":64500}]}"#

    static let bookingNotFound = #"{"error":"Booking not found"}"#

    static let operatorBookings = #"{"bookings":[{"id":1,"status":"paid","total_cents":129000,"created_at":"2027-03-04T10:21:09Z"},{"id":2,"status":"pending","total_cents":64500,"created_at":"2027-03-04T10:33:51Z"}]}"#

    static let plans = #"{"plans":[{"id":1,"name":"Single Pass","description":"One contact, any station.","price_cents":64500,"stations":"any"},{"id":2,"name":"Polar Bundle","description":"Ten passes over Ystmark.","price_cents":580000,"stations":"YST"},{"id":3,"name":"Salt Flat Hour","description":"Low mask, dry air.","price_cents":71000,"stations":"SQT"},{"id":4,"name":"Equator Week","description":"Low-inclination orbits.","price_cents":220000,"stations":"LAN"},{"id":5,"name":"Storm Spare","description":"A backup pass when the Arctic stows.","price_cents":49000,"stations":"CWR"}]}"#

    static let billingExport = #"{"integration":"Ledger Export","bookings":[{"id":1,"status":"paid","total_cents":129000,"created_at":"2027-03-04T10:21:09Z"},{"id":2,"status":"paid","total_cents":64500,"created_at":"2027-03-04T10:33:51Z"}]}"#

    static let keyChallenge = #"{"challengeId":"8a2c5e1f-7b3d-4f96-a0e4-6d9b1c3f5a72","publicKey":{"challenge":"t7-Kp2Qx9sR4mW1vY6uB3nL8cF5hJ0dZ2gA4eC7","timeout":60000,"rpId":"localhost","allowCredentials":[],"userVerification":"required"}}"#
}
