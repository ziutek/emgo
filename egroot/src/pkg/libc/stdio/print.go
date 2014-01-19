package unistd

func Print(s string) // defined in print.c

func Println(s string) {
	Print(s)
	Print("\n")
}
