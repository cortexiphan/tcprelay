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
	addr, _ := net.ResolveTCPAddr("tcp", ":8080")
	l, _ := net.ListenTCP("tcp", addr)
	for {
		conn, _ := l.AcceptTCP()
		f, _ := os.Create("dst.txt")
		w := bufio.NewWriter(f)
		n, err := io.Copy(w, conn)
		if err != nil {
			fmt.Printf("copy fail, err:%v\n", err)
			conn.Close()
			f.Close()
			continue
		}
		conn.Close()
		fmt.Printf("total copied bytes count:%d\n", n)
		data, _ := os.ReadFile("dst.txt")
		sum := md5.Sum(data)
		fmt.Printf("md5:%x\n", sum)
	}
}
