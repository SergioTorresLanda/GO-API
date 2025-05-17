package main

import (
	"fmt"
	"math/rand"
	"time"
)

//CHANNELS
//1.HOLD DATA
//2.THREAD SAFE (Avoid data races when writing reading)
//3.LISTEN FOR DATA
///contains an underlying array
func mainCN(){
	var c = make(chan int, 5)  //buffer (multiple values)
	go process(c)
	///Dead lock error ocurrs if this keeps waiting for updates on the channel
	for i:= range c {
		fmt.Println(i)
		time.Sleep(time.Second*1) //Some work..
	}
	//fmt.Println(<-c)  //value gets popped out the channel into variable i 
}

func process(c chan int){
	defer close(c) //To AVOID DEADLOCK ERROR !! //defer = do this before the functions exits...
	for i:=0; i<5; i++{
		c <- i
	}
	fmt.Println("exiting process") 
	//Will exit early without having to wait for the main func to make room in the channel if is buffer
	//If its unbuffer will exit at the end because it has to wait for the main func to pop out the value on the channel.
	//close(c) //also can be called here
}

//MORE REAL EXAMPLE
var MAX_CHICKEN_PRICE float32 = 5
var MAX_TOFU_PRICE float32 = 3

func mainChan2(){
	var chickenChannel = make(chan string)
	var tofuChannel = make(chan string)
	var websites = []string{"walmart.com", "costco.com", "wholefoods.com"}
	for i:= range websites{
		go checkChickenPrices(websites[i], chickenChannel)
		go checkTofuPrices(websites[i], tofuChannel)
	}
	sendmessage(chickenChannel, tofuChannel)
}

func checkChickenPrices(website string, chickenChannel chan string){
	for {
		time.Sleep(time.Second*1)
		var chicken_price = rand.Float32()*20
		if chicken_price<MAX_TOFU_PRICE{
			chickenChannel <- website
			break
		}
	}
}

func checkTofuPrices(website string, tofuChannel chan string){
	for {
		time.Sleep(time.Second*1)
		var tofu_price = rand.Float32()*20
		if tofu_price<MAX_TOFU_PRICE{
			tofuChannel <- website
			break
		}
	}
}

func sendmessage(chickenChannel chan string, tofuChannel chan string){
	select{ //SWITCH STATEMENT FOR CHANNELS
		case website := <-chickenChannel:
			fmt.Printf("\nTEmail Sent: Found deal on chicken at %v", website)
		case website := <-tofuChannel:
			fmt.Printf("\nTEmail Sent: Found deal on tofu at %v", website)
	}
}

