package main

func main() {
	println(greeting("World"))
}

func greeting(name string) string {
	return "Hello, " + name + "!"
}
