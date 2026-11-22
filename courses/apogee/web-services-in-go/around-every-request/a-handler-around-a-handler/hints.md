`http.HandlerFunc` turns a function with a handler's two parameters into an `http.Handler`, so `noStore` can return one built around a function literal that receives `w` and `r`.

Inside that function, set the header through `w.Header().Set`, then pass the request on with `next.ServeHTTP`, handing it the same `w` and `r`.

If the header is missing even though you set it, check the order: a header set after `next.ServeHTTP` returns comes too late, because the answer was already written inside that call.
