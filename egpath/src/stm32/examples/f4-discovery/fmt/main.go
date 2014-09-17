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

	i := fmt.Int64(15)
	i.Format(con, 10, 5)
	end.Format(con)
	i.Format(con, 10, -5)
	end.Format(con)
	i.Format(con, -10, 5)
	end.Format(con)
	i.Format(con, -10, -5)
	end.Format(con)

	s := fmt.Str("abcd")
	s.Format(con, 11)
	end.Format(con)
	s.Format(con, -11)
	end.Format(con)

	fmt.Fprint(con, s, fmt.T, i, fmt.N)
}
