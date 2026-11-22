Hash the whole key, the customer ID and the dot included, with `sha256.Sum256([]byte(key))`: `KeyHash` was made from the whole key.

`hex.EncodeToString` needs a slice, so pass it `sum[:]`. It writes lowercase digits, the same as `KeyHash`.

If Furrowlight's key is still refused, check that you compare the 64 hex digits, not the raw 32 bytes, with `KeyHash`, both as `[]byte`, in the same `subtle.ConstantTimeCompare` call as before.
