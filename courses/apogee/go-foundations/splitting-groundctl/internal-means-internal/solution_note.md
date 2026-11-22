`main` sits outside `uplink/`, so the go command refuses its import of `uplink/internal/keys`; `uplink` itself may import it, and its exported `KeyID` hands out only the ID, never the secret.
