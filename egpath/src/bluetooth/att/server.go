package att

type Handler interface {
	ServeATT(w *ResponseWriter, r *Request)
}
