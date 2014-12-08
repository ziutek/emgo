package main

import (
	"fmt"
	"stm32/f4/setup"
)

func init() {
	setup.Performance168(8)
	initConsole()
}

func main() {
	con.WriteString("fmt test:\n")
	end := fmt.Str("$\n")

	s := fmt.Str("abcd")
	i := fmt.Int64(15)

	i.Format(con, 5, 10)
	end.Format(con)
	i.Format(con, -5, 10)
	end.Format(con)
	i.Format(con, 5, -10)
	end.Format(con)
	i.Format(con, -5, -10)
	end.Format(con)

	s.Format(con, 11)
	end.Format(con)
	s.Format(con, -11)
	end.Format(con)

	fmt.Fprint(con, s, fmt.T, i, fmt.N)
	fmt.Fprintf(con, "%v:\t:%v\n", s, i)

	con.WriteString("\nSome error messages:\n\n")

	fmt.Fprintf(con, "%v %v %v %v\n", s, i)
	fmt.Fprintf(con, "%v\n", s, s, i)
	fmt.Fprint(con, fmt.N)
}
