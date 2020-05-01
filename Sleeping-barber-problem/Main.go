package main

import (
	"fmt"
	"sync"
	"time"
)

type Customer struct {
	name string
}
type Barber struct {
	name string
	sync.Mutex
	state    string
	customer *Customer
}

var waitingGroup *sync.WaitGroup

func (c *Customer) String() string {
	return fmt.Sprintf("%p", c)[7:]
}

func BarberUncle() (b *Barber) {
	return &Barber{
		name:  "barber uncle",
		state: "sleeping",
	}
}

func barber(barber *Barber, waitingRoom chan *Customer, newCustomer chan *Customer) {
	for {
		barber.Lock()
		defer barber.Unlock()
		barber.state = "checking"
		barber.customer = nil

		fmt.Printf("Checking waiting room, has  %d customers\n", len(waitingRoom))
		time.Sleep(time.Microsecond * 100)
		select {
		case customer := <-waitingRoom:
			HairCut(customer, barber)
			barber.Unlock()
		default:
			fmt.Printf("Waiting room is empty, %s goes to sleep\n", barber.name)
			barber.state = "sleeping"
			barber.customer = nil
			barber.Unlock()
			customer := <-newCustomer
			barber.Lock()
			fmt.Printf("%s is woken up by %s\n", barber.name, customer)
			HairCut(customer, barber)
			barber.Unlock()
		}
	}
}

func HairCut(customer *Customer, barber *Barber) {
	barber.state = "cutting"
	barber.customer = customer
	barber.Unlock()
	fmt.Printf("Doing haircut of a %s\n", customer)
	time.Sleep(time.Microsecond * 100)
	barber.Lock()
	waitingGroup.Done()
	fmt.Printf("%s haircut is done\n", customer)
	barber.customer = nil
}

func customer(customer *Customer, barber *Barber, waitingRoom chan *Customer, newCustomer chan *Customer) {
	time.Sleep(time.Microsecond * 	50) //arrival
	barber.Lock()
	fmt.Printf("new customer %s comes and sees %s barber, serving customer %s. WaitingRoom curerntly has %d customers\n",
		customer,
		barber.state,
		barber.customer,
		len(waitingRoom))
	switch barber.state {
	case "sleeping":
		select {
		case newCustomer <- customer:
		default:
			select {
			case waitingRoom <- customer:
			default:
				waitingGroup.Done()
			}

		}
	case "cutting":
		select {
		case waitingRoom <- customer:
		default:
			waitingGroup.Done()
		}
	case "checking":
		panic("Customer should not checking for waiting room when barber is checking for customer")
	}
	barber.Unlock()
}



func main() {
	barberUncle := BarberUncle()
	waitingRoom := make(chan *Customer, 5)
	servingCustomer := make(chan *Customer, 1)
	go barber(barberUncle, waitingRoom, servingCustomer)

	time.Sleep(time.Millisecond * 100)

	waitingGroup = new(sync.WaitGroup)
	n := 10
	waitingGroup.Add(10)

	for i := 0; i < n; i++ {
		time.Sleep(time.Millisecond * 50)
		newCustomer := new(Customer)
		go customer(newCustomer, barberUncle, waitingRoom, servingCustomer)
	}
	waitingGroup.Wait()
	fmt.Println("No more customers for the day, closing No1 salon for the quality haircut")
}
