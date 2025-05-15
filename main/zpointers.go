package main 
import (
	"fmt"
)

func mainZ(){
//POINTERS (variables that stores memory location)

var p *int32 //nil
var p2 *int32 = new(int32)  //0xc00009000c
var i int32 = 55 // (0x1b08)
*p2=10 //make sure p is not nil before accesing assing value..
fmt.Println(p)
fmt.Println(p2) //pointer-memory loc.
fmt.Println(*p2) //value

p3 := &i //reference memory address of i  
fmt.Println(p3) //pointer  (0xc000104020)
fmt.Println(*p3) //value (55)
*p3=44
fmt.Println(i) //value of i changes !! (44)

var thing1 = [5]float64{1,2,3,4,5}
result := square(&thing1)
fmt.Printf("Result is: %v", result)
fmt.Printf("thing1 is: %v", thing1)


//CONLUSION FOR LARGE PARAMETERS USE POINTERS TO SAVE MEMORY!! 
}

func square(thing2 *[5]float64) [5]float64{
	fmt.Printf("The memory alloc is: %p", thing2) 
	for i := range thing2 {
		thing2[i] = thing2[i]*thing2[i] //also changes thing1 !!
	}
	return *thing2
}