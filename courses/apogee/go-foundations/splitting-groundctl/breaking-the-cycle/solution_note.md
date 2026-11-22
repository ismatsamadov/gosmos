`severity` imports nothing from groundctl, so `telemetry` and `alert` can both import it. `alert` still imports `telemetry`, but no import leads back to `alert`, so the cycle is gone.
