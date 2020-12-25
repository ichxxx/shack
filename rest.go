package shack

type Resp struct {
	Code  string  `json:"code,omitempty"`
	Msg   string  `json:"msg,omitempty"`
	Error string  `json:"error,omitempty"`
	Data  map[string]interface{}
}


func(r *Resp) Success() {

}
