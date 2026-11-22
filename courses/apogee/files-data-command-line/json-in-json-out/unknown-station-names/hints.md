Declare a `map[string]int` variable and pass its address to `json.Unmarshal`; it makes the map for you and adds every key it finds.

When `Unmarshal` fails, the map may already hold some stations, such as `SQT` before a bad `CWR`. Return nil with the error instead of the map.
