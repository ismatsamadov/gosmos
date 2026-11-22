Each tag gives a field its key in the export, and encoding/json uses the tag in both directions. Untagged fields would still decode, but they would be written back under their Go names.
