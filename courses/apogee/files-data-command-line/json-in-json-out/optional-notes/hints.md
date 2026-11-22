The option goes inside the tag's quotes, after the key name and a comma, with no spaces.

Leave the `confirmed` tag as it is: with `omitempty` there, every booking that is not confirmed yet would lose its false. Then return what `json.Marshal` gives for the confirmation.
