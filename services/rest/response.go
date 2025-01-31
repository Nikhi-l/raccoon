package rest

import (
	"io"
	pb "raccoon/proto"
	"raccoon/serialization"
)

type Response struct {
	*pb.EventResponse
}

func (r *Response) SetCode(code pb.Code) *Response {
	r.Code = code
	return r
}

func (r *Response) SetStatus(status pb.Status) *Response {
	r.Status = status
	return r
}

func (r *Response) SetSentTime(sentTime int64) *Response {
	r.SentTime = sentTime
	return r
}

func (r *Response) SetReason(reason string) *Response {
	r.Reason = reason
	return r
}

func (r *Response) SetDataMap(data map[string]string) *Response {
	r.Data = data
	return r
}

func (r *Response) Write(w io.Writer, s serialization.Serializer) (int, error) {
	b, err := s.Serialize(r)
	if err != nil {
		return 0, err
	}
	return w.Write(b)
}
