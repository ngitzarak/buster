package     main

import "fmt"
import "net"
//import "net/socks"
import "golang.org/x/net/internal/socks"
import "io"
import "os"

func handle(local_conn, socks_conn net.Conn) {

  	go func() {
		defer local_conn.Close()
		defer socks_conn.Close()
		io.Copy(local_conn, socks_conn)
	}()
	go func() {
		defer local_conn.Close()
		defer socks_conn.Close()
		io.Copy(socks_conn, local_conn)
	}()

}
func handleIncoming(ln net.Listener, address string) {
  for {
      conn, err := ln.Accept()
      if err != nil {
          fmt.Println("Couldn't accept connection")
          return
      }

        dialer := socks.NewDialer("tcp", "127.0.0.1:9050")
        socks_conn, err := dialer.Dial("tcp", address)
        if err != nil {
            conn.Close()
            fmt.Println("Couldn't dial")
            continue
        }
        go handle(conn, socks_conn)


  }
}

func main() {
    ln, err := net.Listen("tcp", os.Args[1])
    if err != nil {
        fmt.Println("Couldn't listen")
        return
    }
    go handleIncoming(ln, os.Args[2])
    
    ch := make(chan []byte)

    <- ch
}
