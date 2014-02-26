package unistd

func Print(s string)

func Println(s string) {
	Print(s)
	Print("\n")
}
