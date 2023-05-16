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
	"sync"
	"syscall"
)

type travelWays struct {
	channels map[cityName]chan *Merchant
	mutex    sync.Mutex
}

func (travelWays *travelWays) Store(city cityName, channel chan *Merchant) {
	travelWays.mutex.Lock()
	defer travelWays.mutex.Unlock()
	if travelWays.channels == nil {
		travelWays.channels = make(map[cityName]chan *Merchant)
	}

	travelWays.channels[city] = channel
}

func (travelWays *travelWays) Load(city cityName) (chan *Merchant, bool) {
	travelWays.mutex.Lock()
	defer travelWays.mutex.Unlock()
	if travelWays.channels == nil {
		return nil, false
	}
	ch, ok := travelWays.channels[city]
	return ch, ok
}

func (travelWays *travelWays) Delete(city cityName) {
	travelWays.mutex.Lock()
	defer travelWays.mutex.Unlock()
	if travelWays.channels == nil {
		return
	}
	delete(travelWays.channels, city)
}

func (travelWays *travelWays) Range(f func(cityName, chan *Merchant) bool) {
	travelWays.mutex.Lock()
	defer travelWays.mutex.Unlock()
	if travelWays.channels == nil {
		return
	}
	for k, v := range travelWays.channels {
		if !f(k, v) {
			break
		}
	}
}

// RegisterTravelWay connects cities using channels
func RegisterTravelWay(fromCity *City, toCity *City) {
	channel := make(chan *Merchant, 100)
	toCity.inboundTravelWays.Store(fromCity.name, channel)
	fromCity.outboundTravelWays.Store(toCity.name, channel)
}

type networkedTravelWays struct {
	city   *City
	server net.Listener
}

// setupNetworkedTravelWay will listen for incoming connections and add them to the cities travelWays. It can also connect to another networkTravelWay
func setupNetworkedTravelWay(portNumber int, city *City) *networkedTravelWays {

	// start a TCP server to listen for requests on
	listener, err := net.Listen("tcp", "localhost:"+strconv.Itoa(portNumber))
	for isErrorAddressAlreadyInUse(err) {
		portNumber++
		listener, err = net.Listen("tcp", "localhost:"+strconv.Itoa(portNumber))
	}
	if err != nil && !isErrorAddressAlreadyInUse(err) {
		fmt.Println(err)
		return nil
	}

	travelWays := &networkedTravelWays{
		city:   city,
		server: listener,
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
			go travelWays.handleConnection(connection)
		}
	}()

	return travelWays
}

func (travelWays *networkedTravelWays) requestConnection(address string) {
	fmt.Printf("Requesting connection: %s...\n", address)

	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	go travelWays.handleConnection(connection)
}

// blocking, must be handled as a new routine
func (travelWays *networkedTravelWays) handleConnection(connection net.Conn) {
	done := make(chan bool)
	outboundChannel := make(chan *Merchant, 100)
	inboundChannel := make(chan *Merchant, 100)

	// send our city's name. Blocking, so run in separate routine
	go func() {
		_, err := connection.Write([]byte(travelWays.city.name))
		if err != nil {
			fmt.Println(err)
			done <- true
		}
	}()

	// get the city's name
	initializationPacketBytes := make([]byte, 1024)
	n, err := connection.Read(initializationPacketBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	remoteCityName := cityName(initializationPacketBytes[:n])

	// make sure we aren't already connected to the city
	if _, alreadyExist := travelWays.city.outboundTravelWays.Load(remoteCityName); alreadyExist {
		fmt.Printf("error: travelWay from %s to %s already exists in outbound travelWays\n", remoteCityName, travelWays.city.name)
		return
	}
	if _, alreadyExist := travelWays.city.inboundTravelWays.Load(remoteCityName); alreadyExist {
		fmt.Printf("error: travelWay from %s to %s already exists in inbound travelWays\n", remoteCityName, travelWays.city.name)
		return
	}

	// add travelWay to city
	travelWays.city.inboundTravelWays.Store(remoteCityName, inboundChannel)
	travelWays.city.outboundTravelWays.Store(remoteCityName, outboundChannel)

	fmt.Printf("Successfully added city %s as a network connection, sending and receiving merchants...\n", remoteCityName)

	go travelWays.handleIncomingMessages(connection, inboundChannel, done)
	go travelWays.handleOutgoingMessages(connection, outboundChannel, done)

	// wait for connection to close
	<-done
	travelWays.city.outboundTravelWays.Delete(remoteCityName)
	travelWays.city.inboundTravelWays.Delete(remoteCityName)

	fmt.Printf("Connection to %s closed\n", remoteCityName)
	connection.Close()
}

// blocking, must be handled as a new routine
func (travelWays *networkedTravelWays) handleIncomingMessages(connection net.Conn, channel chan *Merchant, done chan bool) {
	defer func() {
		done <- true
	}()

	// pass merchants from connection to channel
	reader := bufio.NewReader(connection)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF && err != syscall.EPIPE {
				fmt.Println(err)
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
func (travelWays *networkedTravelWays) handleOutgoingMessages(connection net.Conn, channel chan *Merchant, done chan bool) {
	defer func() {
		done <- true
	}()

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
