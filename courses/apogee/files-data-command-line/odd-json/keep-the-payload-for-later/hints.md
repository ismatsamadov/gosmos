With `any`, `json.Unmarshal` decodes the payload straight away into a map, and its spacing and order are gone. A `json.RawMessage` field keeps the payload's bytes instead.

A report without a payload leaves the field empty, and `"payload":null` fills it with the four bytes `null`. After `Unmarshal` succeeds, check for both and return an error that says the report has no payload.
