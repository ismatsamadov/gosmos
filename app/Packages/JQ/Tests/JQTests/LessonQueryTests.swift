import JQ
import Testing

/// What a lesson expects of a query's output.
enum Expected: Sendable {
    /// One output equal to this: the lessons' `==` and `eq`.
    case equals(JQValue)
    /// One number greater than this: `>` and `gt`.
    case greaterThan(Int)
    /// One number at least this: `>=`.
    case atLeast(Int)
    /// One string containing this: `contains`.
    case contains(String)
    /// One string without this: `not_contains`.
    case lacks(String)
}

/// A lesson's query with a realistic input and the answer the lesson expects.
struct LessonCase: Sendable, CustomTestStringConvertible {
    let query: String
    let mode: JQInputMode
    let input: String
    let expected: Expected

    var testDescription: String {
        query
    }

    static func json(_ query: String, _ input: String, _ expected: Expected) -> LessonCase {
        LessonCase(query: query, mode: .json, input: input, expected: expected)
    }

    static func jsonl(_ query: String, _ input: String, _ expected: Expected) -> LessonCase {
        LessonCase(query: query, mode: .jsonl, input: input, expected: expected)
    }

    func verify() throws {
        let outputs = try JQProgram(query).run(JQValue(json: input, mode: mode))
        switch expected {
        case .equals(let value):
            #expect(outputs == [value])
        case .greaterThan(let bound):
            let number = try #require(Self.number(outputs))
            #expect(number > Double(bound))
        case .atLeast(let bound):
            let number = try #require(Self.number(outputs))
            #expect(number >= Double(bound))
        case .contains(let text):
            let string = try #require(Self.string(outputs))
            #expect(string.contains(text))
        case .lacks(let text):
            let string = try #require(Self.string(outputs))
            #expect(!string.contains(text))
        }
    }

    private static func number(_ outputs: [JQValue]) -> Double? {
        guard outputs.count == 1, case .number(let number) = outputs[0] else {
            return nil
        }
        return number.doubleValue
    }

    private static func string(_ outputs: [JQValue]) -> String? {
        guard outputs.count == 1, case .string(let string) = outputs[0] else {
            return nil
        }
        return string
    }
}

extension JQTests {
    /// Lesson-style jq queries, written for these tests, in every place a lesson uses one:
    /// command-line output tests, HTTP captures and HTTP tests.
    @Suite("Lesson queries")
    struct LessonQueries {
        @Test("Check program reports", arguments: Lessons.checkReports)
        func checkReport(_ query: String) throws {
            let program = try JQProgram(query)
            let fields = query.matches(of: /\.([A-Za-z][A-Za-z0-9]*)/).map { String($0.1) }
            try #require(!fields.isEmpty)
            // The check programs print one boolean per property, true when it holds.
            func report(failing: String?) throws -> JQValue {
                try JQValue(json: "{" + fields.map { #""\#($0)":\#($0 != failing)"# }.joined(separator: ",") + "}")
            }
            #expect(try program.run(report(failing: nil)) == [true])
            for field in fields {
                #expect(try program.run(report(failing: field)) == [false], "with \(field) false")
            }
        }

        @Test("Command-line output tests", arguments: Lessons.commandLine)
        func commandLine(_ lesson: LessonCase) throws {
            try lesson.verify()
        }

        @Test("HTTP response captures", arguments: Lessons.captures)
        func capture(_ lesson: LessonCase) throws {
            try lesson.verify()
        }

        @Test("HTTP JSON tests", arguments: Lessons.httpTests)
        func httpTest(_ lesson: LessonCase) throws {
            try lesson.verify()
        }
    }
}

enum Lessons {
    /// Queries over the reports of check programs, which all pass when every named property
    /// is true.
    static let checkReports: [String] = [
        #".aosBeforeLos and .losAfterAos and .durationMatchesTable"#,
        #".signatureChecked and .replayRefused and .staleCommandRefused and .keyRotationHonoured"#,
        #".bookingStored"#,
        #".maskRespected and .horizonLoaded and .lowPassesSkipped and .overlapsRefused"#,
        #".windStowHonoured and .warningsLogged"#,
        #".utcEverywhere"#,
        #".dopplerCorrected and .beaconLocked and .offsetLogged"#,
        #".framesComplete and .checksumsMatch and .deliveredBeforeDeadline and .operatorNotified"#,
        #".queueBounded and .retriesSpaced and .deadLettersKept"#,
        #".keysHashed and .rawKeysAbsent and .rotationLogged"#,
        #".boardSorted and .nextPassFirst and .pastPassesHidden"#,
        #".backhaulFallback and .fiberRestoredDetected and .noFramesLost and .operatorsWarned and .stowLogged"#,
    ]

    /// Other command-line output tests: Go's environment, structured reports, and JSON-lines
    /// logs.
    static let commandLine: [LessonCase] = [
        .json(#".GOVERSION | [scan("[0-9]+") | tonumber] as [$major, $minor, $patch] | $major > 1 or $minor > 27 or ($minor == 27 and $patch >= 2)"#, Fixture.goEnv("go1.27.3"), .equals(true)),
        .json(#".revoked == 2 and (.oldKeyAccepted | not) and (.graceKeyAccepted | not) and .newKeyAccepted and .revokedAtMatchesClock and .secondRevocationNoChange"#, Fixture.keyRevocationReport, .equals(true)),
        .json(#".aosBeforeLos and .durationSeconds > 600 and .maxElevationAboveMask and .withinBooking"#, Fixture.passTimingReport, .equals(true)),
        .json(#".bookingStatuses == [201, 201, 201, 201, 429] and (.stationStatuses[0:30] | all(. == 200)) and .stationStatuses[30] == 429 and (.commandStatuses[:12] | all(. == 202)) and .commandStatuses[-1] == 429 and .retryAfterOnEveryRefusal and .refusalBodiesMatch"#, Fixture.burstReport, .equals(true)),
        .json(#".pendingChallenge"#, #"{"pendingChallenge":true,"challengeId":17}"#, .equals(true)),
        .json(#".pendingChallenge"#, #"{"pendingChallenge":false}"#, .equals(false)),

        .jsonl(#".[] | select(.msg == "uplink refused") | .error.stack | length"#, Fixture.schedulerLog, .greaterThan(0)),
        .jsonl(#".[] | select(.msg == "listing downlinks failed") | .error.path"#, Fixture.listFailedLog, .equals("archive/CORRUPT")),
        .jsonl(#".[] | select(.msg == "listing downlinks failed") | .errors | map(.path) | sort | join(", ")"#, Fixture.joinedErrorsLog, .equals("archive/TORN-A, archive/TORN-B")),
        .jsonl(#".[] | select(.msg == "listing downlinks failed") | .errors.error_1 != null"#, Fixture.joinedErrorsLog, .equals(true)),
        .jsonl(#".[] | select(.msg == "listing downlinks failed") | .errors.error_2 != null"#, Fixture.joinedErrorsLog, .equals(true)),
        .jsonl(#"[.[] | select(.msg == "pass planned")] | length"#, Fixture.schedulerLog, .equals(1)),
        .jsonl(#"[.[] | select(.msg == "uplink refused")] | any(.operator == "furrowlight")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"[.[] | select(.msg | contains("pass skipped") or contains("pass dropped"))] | length"#, Fixture.schedulerLog, .equals(0)),
        .jsonl(#"all(.[]; .station | test("^[A-Z]{3}$"))"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"[.[] | .trace_id // empty] | unique | length"#, Fixture.schedulerLog, .equals(6)),
        .jsonl(#"[.[] | select(.msg == "handled call") | .path | test("^/v1/")] | map(select(.)) | length"#, Fixture.schedulerLog, .equals(5)),
        .jsonl(#".[] | select(.level == "WARN") | .wind_kmh > 50"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | .msg; contains("archive gone"))"#, Fixture.missingFileLog, .equals(false)),
        .jsonl(#"any(.[] | .msg; contains("read archive/CORRUPT"))"#, Fixture.listFailedLog, .equals(false)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/"); .client == "10.20.0.x")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/v1/bookings" and .status == 201); .operator == "[REDACTED]")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/v1/bookings" and .status == 201); .operator == null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/v1/commands" and .status == 403); .error.message == "command refused")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/v1/commands" and .status == 403); .operator != null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .status == 500); .error == null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .status == 500); .error.stack == null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[] | select(.msg == "handled call" and .path == "/v1/stations/XXX" and .status == 404); .error.message == "no such station")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call"); .duration_us > 0)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call"); .request_bytes >= 0)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call"); .response_bytes > 0)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "handled call"); .status == 500)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "pass planned"); .aos < .los)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "uplink refused") | .error; contains("uplink.(*Sender).Verify"))"#, Fixture.textStackLog, .equals(true)),
        .jsonl(#"any(.[] | select(.msg == "downlink file missing") | .error; contains("p-0303-sqt-04"))"#, Fixture.missingFileLog, .equals(true)),
        .jsonl(#"any(.[]; .event == "console_opened" and .operatorId != null and .email == "ops@furrowlight.example" and .expiresAt != null and .sessionId == null and .sessionToken == null)"#, Fixture.auditLog, .equals(true)),
        .jsonl(#"any(.[]; .event == "booking_paid" and .email == "[REDACTED]" and .billingName == "[REDACTED]" and .billingAddress == "[REDACTED]" and .billingCity == "[REDACTED]" and .billingCountry == "[REDACTED]" and .billingPostcode == "[REDACTED]" and (.operatorId | type) == "number" and (.bookingId | type) == "number" and (.totalCents | type) == "number")"#, Fixture.auditLog, .equals(true)),
        .jsonl(#"any(.[]; .event == "sign_in" and .email == "[REDACTED]" and (.operatorId | type) == "number" and .success == true)"#, Fixture.auditLog, .equals(true)),
        .jsonl(#"any(.[]; .event == "key_reset_request" and .email == "ops@furrowlight.example" and .success == true and .resetLink != null)"#, Fixture.auditLog, .equals(true)),
        .jsonl(#"any(.[]; .hostname == null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/" and .trace_id != null and .trace_id != "")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/" and .trace_id == "tr-first-light")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/" and .operator != null)"#, Fixture.schedulerLog, .equals(false)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/v1/passes" and .status == 200)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "handled call" and .status == 403)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/v1/bookings" and .status == 500)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "handled call" and .path == "/v1/stations/XXX" and .status == 404)"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "uplink refused" and .build_time != null and .build_time != "unknown")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "uplink refused" and .build != null and .build != "unknown")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "uplink refused")"#, Fixture.schedulerLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "uplink refused")"#, Fixture.textStackLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "listing downlinks failed")"#, Fixture.listFailedLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "listing downlinks failed")"#, Fixture.joinedErrorsLog, .equals(true)),
        .jsonl(#"any(.[]; .msg == "downlink file missing")"#, Fixture.missingFileLog, .equals(true)),
        .jsonl(#"length"#, Fixture.schedulerLog, .atLeast(3)),
    ]

    /// Paths of values captured from a response for later steps.
    static let captures: [LessonCase] = [
        .json(#".[0].id"#, Fixture.downlinksOnEdge, .equals(.string(Fixture.fileID))),
        .json(#".[0].preview_url"#, Fixture.downlinksServedByAPI, .equals(.string("http://localhost:7340/v1/previews/\(Fixture.fileID)"))),
        .json(#".[0].file_url"#, Fixture.downlinksInBuckets, .equals("https://apogee-downlinks.storage.example/frames/aB3nQ8mK4xR6bW1cY9sL2dF7.bin")),
        .json(#".[1].file_url"#, Fixture.downlinksInBuckets, .equals("https://apogee-downlinks.storage.example/telemetry/hJ5gP0aZ8kU2eN6rX4wB1qM9.bin")),
        .json(#".challengeId"#, Fixture.keyChallenge, .equals("8a2c5e1f-7b3d-4f96-a0e4-6d9b1c3f5a72")),
        .json(#".id"#, Fixture.createdOperator, .equals(.string(Fixture.operatorID))),
        .json(#".publicKey.challenge"#, Fixture.keyChallenge, .equals("t7-Kp2Qx9sR4mW1vY6uB3nL8cF5hJ0dZ2gA4eC7")),
        .json(#".publicKey.rpId"#, Fixture.keyChallenge, .equals("localhost")),
        .json(#".refresh_token"#, Fixture.signIn, .equals("4d8a1f6c3b9e2d7a0f5c8b1e4d7a2f9c6b3e0d8a5f2c9b6e3d0a7f4c1b8e5d2a")),
        .json(#".token"#, Fixture.signIn, .equals("eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhcG9nZWUtYm9va2luZyIsInN1YiI6ImM3MWQzYTk1IiwiZXhwIjoxODA0MDAwMDAwfQ.Zx4kQ9wT1mB7nR2vY8cL5pH3sJ6dF0gA")),
        .json(#"first(.[] | select(.title == "The Board")).uid"#, Fixture.boardSearch, .equals("apogee-next-passes")),
    ]

    /// Paths of HTTP steps' JSON tests.
    static let httpTests: [LessonCase] = [
        .json(#"((.booking | keys | sort) == ["created_at", "id", "status", "total_cents"]) and (.booking.id == 1) and all(.lines[];\#n  ((keys | sort) == ["minutes", "pass_id", "price_cents", "station"]))"#, Fixture.paidBooking, .equals(true)),
        .json(#". | length"#, Fixture.boardSearch, .greaterThan(0)),
        .json(#".[0].satellite"#, Fixture.passesByTime, .equals("pipit-1")),
        .json(#".[0].satellite"#, Fixture.passesByElevation, .equals("cirrine-1")),
        .json(#".[0].satellite"#, Fixture.passes(["thistle", "brine-3"]), .equals("thistle")),
        .json(#".[0].description"#, Fixture.downlinksOnEdge, .contains("over the Ystmark fjords")),
        .json(#".[0].satellite"#, Fixture.routeBindings, .equals("pipit-1")),
        .json(#".[0].key"#, Fixture.routeBindings, .equals("pipit")),
        .json(#".[0].preview_url"#, Fixture.downlinksServedByAPI, .contains("/v1/previews")),
        .json(#".[0].preview_url"#, Fixture.downlinksInBuckets, .contains("/archive/")),
        .json(#".[0].preview_url"#, Fixture.downlinksWithInlinePreviews, .contains("data:image/")),
        .json(#".[0].preview_url"#, Fixture.downlinksInBuckets, .lacks(Fixture.fileID)),
        .json(#".[0].title"#, Fixture.downlinksOnEdge, .contains("Furrow-A swath")),
        .json(#".[0].file_url"#, Fixture.downlinksInBuckets, .contains("/frames/")),
        .json(#".[0].file_url"#, Fixture.downlinksInBuckets, .contains("storage.example")),
        .json(#".[0].file_url"#, Fixture.downlinksSigned, .contains("signature=")),
        .json(#".[0].file_url"#, Fixture.downlinksOnEdge, .contains("edge.apogee.example")),
        .json(#".[1].satellite"#, Fixture.passesByTime, .equals("furrow-a")),
        .json(#".[1].satellite"#, Fixture.passesByElevation, .equals("furrow-b")),
        .json(#".[1].file_url"#, Fixture.downlinksOnEdge, .contains("/telemetry/")),
        .json(#".[2].satellite"#, Fixture.passesByTime, .equals("furrow-b")),
        .json(#".[2].satellite"#, Fixture.passesByElevation, .equals("furrow-a")),
        .json(#".[3].satellite"#, Fixture.passesByTime, .equals("cirrine-1")),
        .json(#".[3].satellite"#, Fixture.passesByElevation, .equals("pipit-1")),
        .json(#".app"#, Fixture.health, .equals("operator-portal")),
        .json(#".accepting"#, Fixture.ystmark, .equals(true)),
        .json(#".accepting"#, Fixture.capeWren, .equals(false)),
        .json(#".satellite"#, Fixture.booking, .equals("pipit-1")),
        .json(#".normalized_callsign"#, #"{"normalized_callsign":"PIPIT-1 / APG-T1"}"#, .equals("PIPIT-1 / APG-T1")),
        .json(#".data.result | length"#, Fixture.gaugeQuery, .greaterThan(0)),
        .json(#".data.result[0].metric.station"#, Fixture.gaugeQuery, .equals("SQT")),
        .json(#".staffed"#, Fixture.ystmark, .equals(true)),
        .json(#".staffed"#, Fixture.capeWren, .equals(false)),
        .json(#".email"#, Fixture.createdOperator, .equals("ops@furrowlight.example")),
        .json(#".email"#, Fixture.updatedOperator, .equals("control@furrowlight.example")),
        .json(#".email"#, Fixture.signIn, .equals("duty@cirrine.example")),
        .json(#".error"#, Fixture.bookingNotFound, .equals("Booking not found")),
        .json(#".mask_deg"#, Fixture.atoll, .equals(8)),
        .json(#".id"#, Fixture.downlink, .equals(.string(Fixture.fileID))),
        .json(#".integration"#, Fixture.billingExport, .equals("Ledger Export")),
        .json(#".integration == "Ledger Export" and (.bookings | length) > 0 and all(.bookings[]; ((keys | sort) == ["created_at", "id", "status", "total_cents"])) and any(.bookings[]; .id == 1 and .total_cents == 129000)"#, Fixture.billingExport, .equals(true)),
        .json(#".priority"#, Fixture.createdOperator, .equals(false)),
        .json(#".priority"#, Fixture.signIn, .equals(true)),
        .json(#".pass_stats.delivered"#, Fixture.ystmark, .greaterThan(1)),
        .json(#".pass_stats.refused"#, Fixture.network, .greaterThan(0)),
        .json(#".pass_stats.planned"#, Fixture.network, .greaterThan(0)),
        .json(#".pass_stats.planned"#, Fixture.saltFlat, .greaterThan(1)),
        .json(#".pass_stats.missed"#, Fixture.stormedYstmark, .greaterThan(100)),
        .json(#".passes_booked"#, Fixture.quietYstmark, .equals(0)),
        .json(#".passes_booked"#, Fixture.ystmark, .greaterThan(2)),
        .json(#".passes_booked"#, Fixture.bookedColumn, .greaterThan(499)),
        .json(#".passes_booked"#, Fixture.saltFlat, .greaterThan(0)),
        .json(#".passes_booked"#, Fixture.atoll, .greaterThan(0)),
        .json(#".ok"#, Fixture.health, .equals(true)),
        .json(#".booking.OperatorEmail"#, Fixture.bookingWithInternalFields, .equals("ops@furrowlight.example")),
        .json(#".booking.ID"#, Fixture.bookingWithInternalFields, .equals(1)),
        .json(#".bookings | length"#, Fixture.billingExport, .greaterThan(0)),
        .json(#".product"#, Fixture.network, .equals("Apogee Ground")),
        .json(#".status"#, Fixture.gaugeQuery, .equals("success")),
        .json(#".kind"#, Fixture.polarRoute, .equals("direct")),
        .json(#".kind"#, Fixture.equatorRoute, .equals("fanout")),
        .json(#".operator_id"#, Fixture.booking, .equals(.string(Fixture.otherOperatorID))),
        .json(#".valid"#, #"{"valid":true}"#, .equals(true)),
        .json(#"all(.bookings[]; ((keys | sort) == ["created_at", "id", "status", "total_cents"]))"#, Fixture.operatorBookings, .equals(true)),
        .json(#"all(.plans[]; ((keys | sort) == ["description", "id", "name", "price_cents", "stations"])) and ([.plans[].name] | sort) == ["Equator Week", "Polar Bundle", "Salt Flat Hour", "Single Pass", "Storm Spare"]"#, Fixture.plans, .equals(true)),
        .json(#"any(.[]; .title == "Uplink Refusal Burst")"#, Fixture.alertRules, .equals(true)),
        .json(#"any(.[]; .title == "Missed Pass Burst")"#, Fixture.alertRules, .equals(true)),
        .json(#"any(.[]; .title == "The Board")"#, Fixture.boardSearch, .equals(true)),
        .json(#"any(.data.groups[]?.rules[]?; .name == "Uplink Refusal Burst" and .state == "firing")"#, Fixture.firingRules, .equals(true)),
        .json(#"any(.data.groups[]?.rules[]?; .name == "Missed Pass Burst" and .state == "firing")"#, Fixture.firingRules, .equals(true)),
        .json(#"first(.[] | select(.title == "Uplink Refusal Burst")).data[0].model.expr"#, Fixture.alertRules, .equals("sum(rate(uplinks_refused_total[5m])) * 60")),
        .json(#"first(.[] | select(.title == "Uplink Refusal Burst")).data[1].model.conditions[0].evaluator.params[0]"#, Fixture.alertRules, .equals(4)),
        .json(#"first(.[] | select(.title == "Uplink Refusal Burst")).group"#, Fixture.alertRules, .equals("Fishbowl")),
        .json(#"first(.[] | select(.title == "Missed Pass Burst")).data[0].model.expr"#, Fixture.alertRules, .equals("sum(increase(passes_missed_total[15m]))")),
        .json(#"first(.[] | select(.title == "Missed Pass Burst")).data[1].model.conditions[0].evaluator.params[0]"#, Fixture.alertRules, .equals(2)),
        .json(#"first(.[] | select(.title == "Missed Pass Burst")).for"#, Fixture.alertRules, .equals("2m")),
        .json(#"first(.[] | select(.title == "Missed Pass Burst")).group"#, Fixture.alertRules, .equals("Fishbowl")),
    ]
}
