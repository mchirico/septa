package logging

import (
	"bytes"
	"fmt"
	"github.com/mchirico/septa/lbytes"
)

func experiment() {
	s := "Hello"

	bufa := lbytes.NewBufferString(s)
	fmt.Fprintf(bufa, "okay")
	fmt.Println(bufa.String())

	buf := bytes.NewBufferString(s)
	fmt.Fprint(buf, ", World!")
	fmt.Println(buf.String())
}
