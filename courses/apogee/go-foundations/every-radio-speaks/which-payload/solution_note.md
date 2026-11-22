`switch item.(type)` compares the type of the value inside `item` with each case, so each payload gets its own route, and `default` catches every other type, such as the stray string.
