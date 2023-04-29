package economy

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

type travelWayInbound interface {
	startCity() cityName
	receiveImmigrant() (bool, *Merchant)
}

type travelWayOutbound interface {
	endCity() cityName
	sendEmigrant(*Merchant)
}

// can be used as either an entering or leaving travelWay
type channeledTravelWay struct {
	city    cityName
	channel chan *Merchant
}

// RegisterChanneledTravelWay connects cities between threads
func RegisterChanneledTravelWay(fromCity *City, toCity *City) {
	channel := make(chan *Merchant, 100)
	toCity.addEnteringTravelWay(&channeledTravelWay{
		city:    fromCity.name,
		channel: channel,
	})
	fromCity.addLeavingTravelWay(&channeledTravelWay{
		city:    toCity.name,
		channel: channel,
	})
}

func (travelWay *channeledTravelWay) endCity() cityName {
	return travelWay.city
}

func (travelWay *channeledTravelWay) startCity() cityName {
	return travelWay.city
}

func (travelWay *channeledTravelWay) receiveImmigrant() (bool, *Merchant) {
	select {
	case merchant := <-travelWay.channel:
		return true, merchant
	default:
		return false, nil
	}
}

// have to be careful, if the receiving city is not popping off merchants, this can become blocking
func (travelWay *channeledTravelWay) sendEmigrant(merchant *Merchant) {
	travelWay.channel <- merchant
}

type NetworkedTravelWays struct {
	city *City
}

// SetupNetworkedTravelWay will listen for incoming connections and add them to the cities travelWays. It can also connect to another networkTravelWay
func SetupNetworkedTravelWay(portNumber int, city *City) *NetworkedTravelWays {
	travelWays := &NetworkedTravelWays{
		city: city,
	}

	// start a TCP server to listen for requests on
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(portNumber))
	for isErrorAddressAlreadyInUse(err) {
		portNumber++
		listener, err = net.Listen("tcp", ":"+strconv.Itoa(portNumber))
	}
	if err != nil && !isErrorAddressAlreadyInUse(err) {
		fmt.Println(err)
		return nil
	}

	// listen for incoming connection requests in the background
	fmt.Printf("%s is listening for tcp connection requests at %s\n", city.name, listener.Addr().String())
	go func() {
		for {
			connection, err := listener.Accept() // blocking call here
			if err != nil {
				fmt.Println(err)
				continue
			}

			// each connection is handled by its own process
			go travelWays.handleIncomingConnection(connection)
		}
	}()

	return travelWays
}

func (travelWays *NetworkedTravelWays) RequestConnection(address string) {
	fmt.Printf("Requesting connection: %s...\n", address)

	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	go travelWays.handleOutgoingConnection(connection)
}

// blocking, must be handled as a new routine
func (travelWays *NetworkedTravelWays) handleIncomingConnection(connection net.Conn) {
	defer connection.Close()

	// get the incoming city's name
	remoteCityNameBytes := make([]byte, 1024)
	n, err := connection.Read(remoteCityNameBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	remoteCityName := cityName(remoteCityNameBytes[:n])

	// send our city's name in response
	_, err = connection.Write([]byte(travelWays.city.name))
	if err != nil {
		fmt.Println(err)
		return
	}

	// add travelWay to city
	channel := make(chan *Merchant, 100)
	travelWays.city.addEnteringTravelWay(&channeledTravelWay{
		city:    remoteCityName,
		channel: channel,
	})

	fmt.Printf("Successfully added city %s as a network connection, accepting merchants...\n", remoteCityName)

	// listen for incoming merchants and pass them into the travelWay
	reader := bufio.NewReader(connection)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			} else {
				// connection broken
				// WARNING: this is accessing a map in another routine, may not be thread safe
				delete(travelWays.city.inboundTravelWays, remoteCityName)
			}
			break
		}

		// Deserialize the merchant object
		merchant := &Merchant{}
		err = json.Unmarshal([]byte(line), merchant)
		if err != nil {
			fmt.Println(err)
			continue
		}

		channel <- merchant
	}
}

// blocking, must be handled as a new routine
func (travelWays *NetworkedTravelWays) handleOutgoingConnection(connection net.Conn) {
	defer connection.Close()

	// send our city's name
	_, err := connection.Write([]byte(travelWays.city.name))
	if err != nil {
		fmt.Println(err)
		return
	}

	// receive the other city's name
	remoteCityNameBytes := make([]byte, 1024)
	n, err := connection.Read(remoteCityNameBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	remoteCityName := cityName(remoteCityNameBytes[:n])

	// add travelWay to city
	channel := make(chan *Merchant, 100)
	travelWays.city.addLeavingTravelWay(&channeledTravelWay{
		city:    remoteCityName,
		channel: channel,
	})

	fmt.Printf("Successfully added city %s as a network connection, pushing merchants...\n", remoteCityName)

	// pass merchants from channel into connection
	writer := bufio.NewWriter(connection)

	for {
		merchant := <-channel

		// Serialize the merchant object
		merchantBytes, err := json.Marshal(merchant)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = writeAndFlush(writer, merchantBytes)
		if err == io.EOF || err == syscall.EPIPE {
			// connection broken, nothing unusual about that
			break
		}
		if err != nil {
			// unexpected error, report it
			fmt.Println(err)
			break
		}
	}

	fmt.Printf("Connection with %s broken, removing travelWay\n", remoteCityName)
	delete(travelWays.city.outboundTravelWays, remoteCityName)

}

func writeAndFlush(writer *bufio.Writer, merchantBytes []byte) error {
	_, err := writer.Write(merchantBytes)
	if err != nil {
		return err
	}

	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}

	return writer.Flush()
}

// helper function from https://stackoverflow.com/a/65865898
func isErrorAddressAlreadyInUse(err error) bool {
	var eOsSyscall *os.SyscallError
	if !errors.As(err, &eOsSyscall) {
		return false
	}
	var errErrno syscall.Errno // doesn't need a "*" (ptr) because it's already a ptr (uintptr)
	if !errors.As(eOsSyscall, &errErrno) {
		return false
	}
	if errErrno == syscall.EADDRINUSE {
		return true
	}
	const WSAEADDRINUSE = 10048
	if runtime.GOOS == "windows" && errErrno == WSAEADDRINUSE {
		return true
	}
	return false
}
