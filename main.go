// main.go

package main

func main() {
	a := App{}
	a.Initialize("root", "root", "nasa")

	a.Run(":8081")
}
