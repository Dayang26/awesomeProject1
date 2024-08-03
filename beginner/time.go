package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	fmt.Println("Current time:", start)
	fmt.Println(strconv.FormatInt(time.Now().Unix(), 10))
	fmt.Println(time.Now().Format(time.RFC850))
	fmt.Println(time.Now().Format(time.RFC1123Z))
	fmt.Println(time.Now().Format("20060102150405"))
	fmt.Println(time.Unix(1401403874, 0).Format("02.01.2006 15:04:05"))
	fmt.Println(time.Unix(1722650048, 0).Format(time.RFC3339))
	fmt.Println(time.Unix(1234567890, 0).Format(time.RFC822))

	elapsed := time.Since(start)

	fmt.Println("time elasped:", elapsed)

}
