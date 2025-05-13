package main 
import ("fmt"
"strings")


func main3(){
	//STRINGS (UTF-8 Encoding) 8bytes
	var myString = "résumé"
	var indexed = myString[0] //indexing the underlying byte array !!!
	fmt.Printf("%v, %T", indexed, indexed) //114, uint8 
	fmt.Printf("the string length is %v", len(myString)) // 8 (number of bytes, no chars!!)
	
	for i, v := range myString{
		fmt.Println(i,v) //will skip indexes of special chars..
	}

	var strSlice = []string{"s","u","x","z"}
	var strBuilder strings.Builder
	for i := range strSlice{
		strBuilder.WriteString(strSlice[i])
	}
	var casStr = strBuilder.String()
	fmt.Println(casStr) //will skip indexes of special chars..


	//RUNES [int32]
	var myRune = []rune("résumé")
	fmt.Printf("the rune length is %v", len(myRune)) // 6 (number of chars !!)
	for i, v := range myRune {
		fmt.Println(i,v) //continuous indexing
	}
	var myRune2 = 'r'
	fmt.Println(myRune2) //114




}