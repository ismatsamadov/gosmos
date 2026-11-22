The minimum is declared by the loop with `:=`, so each retained closure can use its own iteration's value.

Use the captured `minimum` inside the function literal; do not introduce one shared minimum outside the loop.
