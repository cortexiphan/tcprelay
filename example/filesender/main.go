package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	addr, _ := net.ResolveTCPAddr("tcp4", ":8082")
	conn, _ := net.DialTCP("tcp4", nil, addr)
	f, _ := os.Open("src.txt")
	r := bufio.NewReader(f)
	n, err := io.Copy(conn, r)
	if err != nil {
		fmt.Printf("copy fail, err:%v\n", err)
		return
	}
	fmt.Printf("total copied bytes count:%d\n", n)
	data, _ := os.ReadFile("src.txt")
	sum := md5.Sum(data)
	fmt.Printf("md5:%x\n", sum)
}
