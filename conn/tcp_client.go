package conn

import(
	"time"
	"net"
	"bufio"
	"fmt"
	"errors"
)

type TcpClient struct {
	IsConnect bool
	conn net.Conn
	host string

	OnRecv func([]byte)
}

func (c *TcpClient) Init(host string) {
	c.IsConnect = false
	c.conn = nil
	c.OnRecv = nil
	c.host = host

	go c.ConnectDemon()

	go func(){
		for {
			if c.IsConnect == false{
				time.Sleep(100 * time.Millisecond)
				continue
			}

			msg,err := bufio.NewReader(c.conn).ReadBytes('\n')
			if err != nil {
				fmt.Println("recv data error: ",err)
				c.conn.Close()
				c.IsConnect = false
				//break
			}else{
				fmt.Println("recv msg : ",msg)
				if c.OnRecv != nil {
					c.OnRecv(msg)
				}
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (c *TcpClient) Send(data []byte) error {
	if c.IsConnect == false {
		return errors.New("Connection not established")
	}

	_,err := c.conn.Write(data)
	if err != nil {
		fmt.Println("Error to send message because of ", err.Error())
		c.IsConnect = false
		c.conn.Close()
	}

	return nil
}

func (c *TcpClient) RegRecver(recver func([]byte)) {
	c.OnRecv = recver
}

func (c *TcpClient) ConnectDemon() {
	var err error
	for {
		if c.IsConnect == false {
			c.conn,err = net.Dial("tcp",c.host)
			fmt.Print("connect (",c.host)
			if err != nil {
				fmt.Println(") fail")
			}else{
				fmt.Println(") ok")
				defer c.conn.Close()
				c.IsConnect = true
			}
		}
		
		time.Sleep(1 * time.Second)
    }
}

func (c *TcpClient) SendTest() {
	go func(){
		for {
			str := "shit xxx"
			err := c.Send([]byte(str))
			if err != nil {
				fmt.Println("send err: ",err)
			}
			
			time.Sleep(1 * time.Second)
		}
	}()
}