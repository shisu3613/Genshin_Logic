package main

//func main() {
//	slice := make([]int, 2, 4)
//	for i := 0; i < len(slice); i++ {
//		slice[i] = i
//	}
//	fmt.Printf("slice: %v, addr: %p \n", slice, slice)
//	changeSlice(slice)
//	fmt.Printf("slice: %v, addr: %p \n", slice, slice)
//}
//func changeSlice(s []int) { //传递的是数组地址
//	s = append(s, 3)
//	//s = append(s, 4)
//	s[1] = 111
//	fmt.Printf("func s: %v, addr: %p \n", s, s)
//}

import "fmt"

func main() {
	fmt.Println("c")
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		fmt.Println("d")
		if err := recover(); err != nil {
			fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		}
		fmt.Println("e")
	}()

	f()              //开始调用f
	fmt.Println("f") //这里开始下面代码不会再执行
}

func f() {
	fmt.Println("a")
	panic("异常信息")
	fmt.Println("b") //这里开始下面代码不会再执行
	fmt.Println("f")
}
