package cassette

// TrimResponseHeaders removes response headers we don't need in place.
func (ias Interactions) TrimResponseHeaders() {
	for _, ia := range ias {
		for k := range ia.Response.Headers {
			switch k {
			case "Content-Type": // leave
			default:
				delete(ia.Response.Headers, k)
			}
		}
	}
}
