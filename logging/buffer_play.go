package logging

import (
	"fmt"
	"github.com/mchirico/septa/lbytes"
)

func experiment() {
	s := "Hello"

	buf := lbytes.NewBufferString(s)
	fmt.Fprintf(buf, " okay")
	fmt.Println(buf.String())

	buf.Slop("log1")
	buf.Slop("log2")

	buf.Slop("log1")

	a := []byte("one")
	buf.Read(a)
	buf.Slop("log2")
	fmt.Printf("bufa.String()= %v\n", buf.String())
	fmt.Printf("a= %s\n", a)

}
