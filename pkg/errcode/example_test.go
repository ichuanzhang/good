package errcode_test

import (
	"errors"
	"fmt"
	"good/pkg/errcode"
)

var (
	NotFoundError = errcode.Register(404, 10404, "not fount")
)

func ExampleError_Wrap() {
	fmt.Println(NotFoundError.Wrap(errors.New("data not in database")))
	// Output: status=404, code=10404, msg=not fount, err=data not in database
}
