package conn

import(
	"github.com/tarm/serial"
	"fmt"
)


type DevServial struct {
	//IsOpen bool
	port   *serial.Port  
	config serial.Config
}


func (s *DevServial)Open(serail_nam string,baud int) error {
	s.config.Name = serail_nam
	s.config.Baud = baud
	var err error
	s.port, err = serial.OpenPort(&s.config)
	if err != nil {
		fmt.Printf("Open serial error: %v\n",err)
		return err
	}

	return nil
}

func (s *DevServial) Send (data []byte) error {
	_, err := s.port.Write(data)
	if err != nil {
		fmt.Printf("Send data error: %v\n",err)
		return err
	}

	return nil
}

func (s *DevServial)Close() {
//	serial.ClosePort(s.Port)
}