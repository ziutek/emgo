package att

import (
	"bluetooth/l2cap"
)

type Handler interface {
	ServeATT(w *ResponseWriter, r *Request)
}

type Server struct {
	req     Request
	w       ResponseWriter
	handler Handler
}

func NewServer(maxMTU int) *Server {
	srv := new(Server)
	srv.w.buf = make([]byte, 23, maxMTU)
	return srv
}

func (srv *Server) SetHandler(h Handler) {
	srv.handler = h
}

// HandleTransaction reads request and writes response to far. This synchronous
// approach is cheap and absolutely sufficient, when one far is used only by
// single application.
func (srv *Server) HandleTransaction(far *l2cap.BLEFAR, cid int) {
	srv.req.readAndParse(far)
	srv.w.next(far, cid)
	srv.handler.ServeATT(&srv.w, &srv.req)
}
