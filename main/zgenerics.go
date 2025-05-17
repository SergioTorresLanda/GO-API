package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

//GENERICS
//COLAPSE DIFERENT METHODS THAT DO THE SAME THING FOR DIFFERENT TYPES !!!

func main(){
	intSlice := []int{1,2,3}
	fmt.Println(sumSlice(intSlice)) //type infered ! [int]
	float32Slice := []float32{2,4,6,1}
	fmt.Println(sumSlice(float32Slice)) //type infered ! [float32]
	boolSlice := []bool{false,true,false}
	fmt.Println(is3Lenght(boolSlice)) //type infered !  [bool]]
	fmt.Println(is3Lenght(float32Slice)) //type infered !

	var contacts []contactInfo = loadJSON[contactInfo]("./contactInfo.json")
	var purchases []purchaseInfo = loadJSON[purchaseInfo]("./purchaseInfo.json")
	fmt.Println(contacts)
	fmt.Println(purchases)


}

func sumSlice[T int | float32 | float64](slice []T) T{ //USE OF GENERIC TYPE "T"
//(Can't use "any" because not all types are compatible with addition operator)
	var sum T 
	for _, v := range slice{
		sum += v
	}
	return sum
}

func is3Lenght[T any](slice []T) bool{ //Can use "any" type for T
	return len(slice)==3
}

func loadJSON[T contactInfo | purchaseInfo](filePath string) []T{ 
	//here method is just recieving string param and type param to return !
	data, _ := os.ReadFile(filePath)

	var loaded = []T{} //declare an empty slice of generic type !
	json.Unmarshal(data, &loaded) //similar to JSON.decode for Swift !!

	return loaded
}

// MODEL
type contactInfo struct{
	Name string
	Email string
}
type purchaseInfo struct{
	Name string
	Price string
	Amount string
}
//Generics also work with structs, to allow multiple types for props.. !!
type car [T gasEngine | eEngine] struct {
	brand string
	model string
	engine T
}