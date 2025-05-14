package main 
import (
	"fmt"
)
//STRUCTS
type gasEngine struct{
	mpg uint8
	gallons uint8
}

type eEngine struct{
	mpkwh uint8
	kwh uint8
}

type gasEngine2 struct{
	mpg uint8
	gallons uint8
	//ownerInfo owner 
	owner //can define abstract types!
	int //can define abstract built-in types!
}
type owner struct{
	name string
}

func mainS(){
	var myEngine gasEngine = gasEngine{mpg: 25, gallons: 5} //possible
	var myeEngine eEngine = eEngine{25, 5} //possible
	var myEngine2 gasEngine2 = gasEngine2{26, 7, owner{"alfrd"} , 5}
	myEngine.mpg = 39 //change prop

	fmt.Println(myEngine.mpg, myEngine.gallons) //acces props
	fmt.Println(myEngine2)

	//Anonymous struct. (define and initialize in the same location)
	myE := struct {
		mpg uint8
		gallons uint8
	}{
		mpg:25, 
		gallons:5,
	}
	fmt.Println(myE)
	// Using anonymous struct in a slice
	people := []struct {
		Name string
		Age  int
	}{
		{"Jane Doe", 25},
		{"Peter Pan", 16},
	}
	
		fmt.Println(people)
	
		// Using anonymous struct as a map value
		employeeData := map[string]struct {
			Role string
			Salary int
		}{
			"John Doe": {"Manager", 60000},
			"Jane Smith": {"Developer", 50000},
		}
		fmt.Println(employeeData["John Doe"].Role)

	fmt.Println(myEngine.milesLeft())
		canMakeIt(myEngine, 89) //You can make it
		canMakeIt(myeEngine, 200) //You will die in the highway
}

func (e gasEngine) milesLeft() uint8 { //defined method for the struct gasEngine
	return e.gallons*e.mpg
}

func (e eEngine) milesLeft() uint8 { //defined method for the struct eEngine
	return e.mpkwh*e.kwh
}
//INTERFACEE !!!
type engine interface{ //like swift protocol !!! :)
	milesLeft() uint8
}

func canMakeIt(e engine, miles uint8){
	if miles<=e.milesLeft(){
		fmt.Println("You can make it")
	}else{
		fmt.Println("You will die in the highway")
	}
}
