package main

//FUNCTIONS and CONTROL STRUCTURES
import (
	"errors"
	"fmt"
)

func main(){
	name := "Sergio"
	printMe(name)

	var num = 11
	var den = 2
	var result, rem, err = intDivision(num, den)

	switch{
	case err!=nil:
		fmt.Printf(err.Error())
	case rem==0:
		fmt.Printf("The result of the integer division is %v ", result,) //print with format
	default:
		fmt.Printf("The result of the integer division is %v with remainder %v", result, rem) //print with format
	}

	switch rem{
	case 0:
		fmt.Println("the division was exact")
	case 1,2:
		fmt.Println("the division was close") 
	default:
		fmt.Println("the division was close")
	}

}

func printMe(value string) {//Param declaration
	fmt.Println(value)
}

func intDivision(num int, den int) (int, int, error) { //multiple params
	var err error // nil
	if den==0{
		err = errors.New("cannot divide by zero")
		return 0, 0, err
	}
	var result int = num/den 
	var rem int = num%den 
	return result, rem, err
}