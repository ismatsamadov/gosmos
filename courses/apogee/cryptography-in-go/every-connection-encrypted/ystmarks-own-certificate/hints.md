A `tls.Certificate` is a struct: its `Certificate` field is a slice of byte slices, one per certificate of the chain in order, and `PrivateKey` is the key of the first. Every `*x509.Certificate` keeps its DER bytes in `Raw`.

Build one `tls.Certificate` from `leaf.Raw` and `inter.Raw`, in that order, with `key` and `leaf`, and return a `tls.Config` whose `Certificates` slice holds it, with `Time: now` beside it.

If the client reports `certificate signed by unknown authority`, the intermediate is missing from the chain; if it reports `not valid for any names`, the intermediate comes first, where the leaf belongs. If `TestClock` fails, `Time` is still unset.
