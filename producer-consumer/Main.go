package main

import "fmt"

type Consumer struct {
	message *chan int
}

type Producer struct {
	message *chan int
	done    *chan bool
}

func MyConsumer(message *chan int) *Consumer {
	return &Consumer{
		message: message,
	}
}

func MyProducer(message *chan int, done *chan bool) *Producer {
	return &Producer{
		message: message,
		done:    done,
	}
}

func (p *Producer) produce(max int) {
	fmt.Println("Produce started")
	for i := 0; i < max; i++ {
		fmt.Printf("Produced: %d\n", i)
		*p.message <- i
	}

	*p.done <- true
	fmt.Println("Produce finished")
}

func (c *Consumer) consume() {
	fmt.Println("Consume started:")
	for {
		message := <- *c.message
		fmt.Printf("Consumed: %d\n", message)
	}
}

func main() {
	 var messageChannel = make(chan int)
	 var doneNotifyChannel = make(chan bool)

	 go MyProducer(&messageChannel, &doneNotifyChannel).produce(5)
	 go MyConsumer(&messageChannel).consume()

	 <-doneNotifyChannel
	 fmt.Println("finishing program as production is done")
}
