
//GO FUNDAMENTALS
//1. Statically Typed Lenguage
//2. Strong Typed Lenguage
//3. GO is compiled (comes with compiler wich translates code into machine code.)
//4. Producing binary files that can be runned as a standalone program) 
package main 
import "fmt"

func main(){
	const pi = 3.1415 //Can't be changed
	var name = "Sergio \nTorres"
	var name2 = `Sergio 
Torres
Landa` //using back quotes to format.
	name3 := "Sergei"  // avoid var keyword and type. 
	var n1 int = 3
	var n2 = 2 //can avoid explicit type if assigned
	var n3 int //default value = 0 , "", false
	//int8 int16, int32, int64 
	//uint8, uint16, uint32, uint64 (positives)
    //use types accordingly to manage memory
	var floatNum float64 = 5.654 //float32
	var myBool = true

	//multiple var init
	var var1, var2 int = 1, 2
	var3, var4 := 3, 4

	fmt.Println(n3)
	fmt.Println(n1/n2) //return integer roounded down.
	fmt.Println(n1%n2) //return remainder = 1
	fmt.Println(pi)
	fmt.Println(name)
	fmt.Println(name2) 
	fmt.Println(floatNum)
	fmt.Println(len("tr4e@#")) 
	fmt.Println(myBool) 
	fmt.Println(name3) 
	fmt.Println(var1+var2, var3+var4) //conncatenation

	for i:=0; i<1000; i++{
		n1+=1
	}
}


// init package (create go.mod file) 
// ~go mod init main
// RUN FILE
// ~ go run filename.go
// ~ go run .
// CREATE BINARY EXECUTABlE
// ~ go build .
