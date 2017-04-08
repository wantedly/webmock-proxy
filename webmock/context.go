package webmock

type Context struct {
	// RequestBody saves request body as byte buffer.
	// goproxy.ProxyCtx has Req field, so we do not need to save *http.Request itself.
	RequestBody []byte
}
