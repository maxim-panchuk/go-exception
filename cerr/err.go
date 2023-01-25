package customerr

import (
	"fmt"
	"net/http"

	routing "github.com/qiangxue/fasthttp-routing"
)

type CustomErr struct {
	ErrorMessage string
	Context      string
	Err          error
	Code         int
}

func (e *CustomErr) Error() string {
	if e.Code == 0 {
		return fmt.Sprintf("message: %s\ncontext: %s\nerr: %v", e.ErrorMessage, e.Context, e.Err)
	}
	return fmt.Sprintf("message: %s\ncontext: %s\nerr: %v, statusCode: %v", e.ErrorMessage, e.Context, e.Err, e.Code)
}

func (e *CustomErr) Temporary() bool {
	return e.Code == http.StatusServiceUnavailable
}

func Wrap(err error, context, message string, code int) *CustomErr {
	return &CustomErr{
		ErrorMessage: message,
		Context:      context,
		Err:          err,
		Code:         code,
	}
}

func HandleError(handler func(*routing.Context) error) routing.Handler {
	return func(c *routing.Context) error {
		err := handler(c)
		if err != nil {
			cerr, ok := err.(*CustomErr)
			if ok {
				if cerr.Temporary() {
					c.Response.SetBody([]byte("Сервис пока не доступен, повторите попытку позже!"))
				} else {
					c.Response.SetBody([]byte("Внутренняя ошибка сервера!"))
				}
				c.Response.SetStatusCode(cerr.Code)
			} else {
				c.Response.SetBody([]byte(err.Error()))
			}
		}
		return err
	}
}
