package errors

import (
	"fmt"
	"testing"
)

func TestStack(t *testing.T) {
	err:=New("hell0")
	err2 := geterr(err)
	fmt.Println(err2)
	fmt.Println(Cause(err2)==err)
	fmt.Println("starc",Stack(err2),"ddddd")
}

func geterr(err error)error{
	return WithStack(err)
}
