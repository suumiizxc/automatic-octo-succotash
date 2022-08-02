package response

import "encoding/json"

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}

func (r *Response) ConvertByte() []byte {
	bytes, _ := json.MarshalIndent(r, "", "\t")
	return bytes
}
