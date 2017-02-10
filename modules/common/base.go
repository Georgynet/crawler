package common

import "github.com/gin-gonic/gin"

type BaseResponse struct {
    Data		interface{}	`json:"data,omitempty"`
    Errors		[]ErrorData	`json:"errors,omitempty"`
}

type ErrorData struct {
    Code	    int64	`json:"code" msgpack:"code"`
    Message	    string	`json:"message" msgpack:"message"`
}

func NewBaseResponse(data interface{}, errors map[int64]string) *BaseResponse {
    base := new(BaseResponse)
    base.Build(data, errors)
    return base
}

func (base *BaseResponse) Build(data interface{}, errors map[int64]string){
    if len(errors) > 0 {
        for key, err := range errors {
            base.AddError(key, err)
        }
    } else {
        base.Data = data;
    }
}

func (base *BaseResponse) AddError(code int64, message string) {
    base.Errors = append(base.Errors, ErrorData {
        Code: code,
        Message: message,
    })
}

func ErrorJSON(c *gin.Context, code int, err string) {
    c.JSON(code, NewBaseResponse(nil, map[int64]string{int64(code): err}))
    c.Abort()
}