package services

type HelloService struct {
}

var DefaultHelloService = &HelloService{}

func (helloService *HelloService) Hello(num int) string {

	// return strings.Repeat("Hello", num)
	return "Hello"
}
