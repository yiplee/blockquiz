package request

import (
	"encoding/json"
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/yiplee/blockquiz/handler/api/errors"
)

func init() {
	gin.DisableBindValidation()
}

type bindObject struct {
	Object interface{}
}

var (
	ErrInvalidParameters = errors.New(10001, "invalid parameters")
)

func BindJSON(c *gin.Context, v interface{}) error {
	body, err := Body(c)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidParameters)
	}

	if _, err := govalidator.ValidateStruct(bindObject{Object: v}); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidParameters)
	}

	return nil
}

func BindQuery(c *gin.Context, v interface{}) error {
	if err := binding.Query.Bind(c.Request, v); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidParameters)
	}

	if _, err := govalidator.ValidateStruct(bindObject{Object: v}); err != nil {
		return fmt.Errorf("%s: %w", err.Error(), ErrInvalidParameters)
	}

	return nil
}
