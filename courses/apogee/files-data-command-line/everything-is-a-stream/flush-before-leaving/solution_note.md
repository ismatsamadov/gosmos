Checking writes catches failures when the buffer fills; checking Flush catches failures on its final block. The helper leaves ownership of the underlying writer with its caller.
