package main

import (
	"fmt"
	"time"
)

func mainx(){

	//ARRAY: 
	//fixed lenght [3]
	//same type (int32)
	//indexable [0-i]
	//contiguous in memory 
	var intArray [3]int32  //[0,0,0] //12bytes (4*3) 
	intArray[1]=1
	fmt.Println(intArray[0]) //access indexed element i
	fmt.Println(intArray[0:3]) //access all 3 elements 
	//Acces memory location (&)
	fmt.Println(&intArray[0]) //0xc000090020

	//Intitialize #2
	var intArray2 [3]int32 = [3]int32{1,2,3} 
	fmt.Println(intArray2) //[1,2,3]
	//Init #3
	intArray3 := [3]int32{5,6,8}
	fmt.Println(intArray3) //[6,7,8]

	//SLICES:
	//Wrappers around arrays
	var intSlice []int32 = []int32{1,2,3} 
	intSlice = append(intSlice, 4)
	fmt.Println(intSlice) //[1,2,3,4] //[1,2,3,4,*,*]
	fmt.Printf("The lenght is: %v the capacity is: %v \n", len(intSlice), cap(intSlice)) //The lenght is: 4 the capacity is: 6
	//intSlice[4] will throw index out of range err. 
	intSlice2 := make([]int32,3,8)
	intSlice = append(intSlice, intSlice2...)
	fmt.Println(intSlice) //[1,2,3,4] //[1,2,3,4,*,*]
	fmt.Printf("The lenght is: %v the capacity is: %v \n", len(intSlice), cap(intSlice)) //The lenght is: 7 the capacity is: 12 (4+8)

	//MAPS
	myMap := make(map[string]uint16) 
	fmt.Println(myMap) //map[]
	myMap2 := map[string]uint16{"Adam":23, "Sarah":45}
	fmt.Println(myMap2) //map[Adam:23 Sarah:45]
	fmt.Println(myMap2["Adam"]) //23
	fmt.Println(myMap2["Adal"]) //0 
	var age, ok = myMap2["Json"] //second optional param to check if map has key
	fmt.Println(age, ok) //0 , false 
	delete(myMap2,"Adam")
	fmt.Println(myMap2) //map[Adam:23]
	myMap2["Julia"]=25
	//LOOPS
	for i:=0; i<10; i++{ //init, condition, post
		fmt.Println("hello") 
	}
	//same as :
	var j int = 0
	for j<10{
		fmt.Println("hello") 
		j = j+1
	}
	for name:= range myMap2{
		fmt.Printf("Name: %v\n",name) //Name: Julia
									  //Name: Sarah
	}
	for i, v := range intSlice{
		fmt.Printf("%v, %v\n",i,v)									
	}
	
}

func main(){
	var n int = 1000000
	var testSlice = []int{}
	var testSliceFast = make([]int, 0, n) //>3x times faster

	fmt.Printf(" Total time without preallocation: %v, \n", timeLoop(testSlice, n))	// 26.178714ms
	fmt.Printf(" Total time with preallocation: %v, \n", timeLoop(testSliceFast, n)) // 7.120449ms, 						

}

func timeLoop(slice []int, n int) time.Duration{
var t0 = time.Now()
for len(slice)<n{
	slice = append(slice, 1)
}
return time.Since(t0)
}