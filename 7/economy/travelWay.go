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
	"strings"
	"syscall"
)

// RegisterTravelWay connects cities using channels
func RegisterTravelWay(fromCity *City, toCity *City) {
	channel := make(chan *Merchant, 100)
	toCity.inboundTravelWays[fromCity.name] = channel
	fromCity.outboundTravelWays[toCity.name] = channel
}

type networkedTravelWays struct {
	city   *City
	server net.Listener
}

type TravelWayProtocolType string

const (
	UNIDIRECTIONAL TravelWayProtocolType = "unidirectional"
	BIDIRECTIONAL  TravelWayProtocolType = "bidirectional"
)

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
			go travelWays.handleIncomingConnection(connection)
		}
	}()

	return travelWays
}

func (travelWays *networkedTravelWays) requestConnection(address string, connectionType TravelWayProtocolType) {
	fmt.Printf("Requesting connection: %s...\n", address)

	connection, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err)
		return
	}

	go travelWays.handleOutgoingConnection(connection, connectionType)
}

// blocking, must be handled as a new routine
func (travelWays *networkedTravelWays) handleIncomingConnection(connection net.Conn) {
	defer connection.Close()

	// check if this is unidirectional or bidirectional and get the city's name
	initializationPacketBytes := make([]byte, 1024)
	n, err := connection.Read(initializationPacketBytes)
	if err != nil {
		fmt.Println(err)
		return
	}
	initializationPacket := string(initializationPacketBytes[:n])
	packetValues := strings.Split(initializationPacket, " ")
	var remoteCityName cityName
	if packetValues[0] == string(BIDIRECTIONAL) {
		remoteListenerIP := packetValues[1]
		travelWays.requestConnection(remoteListenerIP, UNIDIRECTIONAL)
		remoteCityName = cityName(packetValues[2])
	} else if packetValues[0] == string(UNIDIRECTIONAL) {
		remoteCityName = cityName(packetValues[1])
	} else {
		fmt.Printf("error: travelWay protocol requires either \"bidirectional ip:port cityName\" or \"unidirectional cityName\", got %s\n", initializationPacket)
		return
	}

	// send our city's name
	_, err = connection.Write([]byte(travelWays.city.name))
	if err != nil {
		fmt.Println(err)
		return
	}

	// make sure we aren't already connected to the city
	if _, alreadyExist := travelWays.city.inboundTravelWays[remoteCityName]; alreadyExist {
		fmt.Printf("error: travelWay from %s to %s already exists\n", remoteCityName, travelWays.city.name)
		return
	}

	// add travelWay to city
	channel := make(chan *Merchant, 100)
	travelWays.city.inboundTravelWays[remoteCityName] = channel

	fmt.Printf("Successfully added city %s as a network connection, accepting merchants...\n", remoteCityName)

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

	fmt.Printf("Connection with %s broken, removing travelWay\n", remoteCityName)
	delete(travelWays.city.inboundTravelWays, remoteCityName)
}

// blocking, must be handled as a new routine
func (travelWays *networkedTravelWays) handleOutgoingConnection(connection net.Conn, connectionType TravelWayProtocolType) {
	defer connection.Close()

	// send connection type
	initializationPacket := ""
	if connectionType == BIDIRECTIONAL {
		initializationPacket = string(BIDIRECTIONAL) + " " + travelWays.server.Addr().String()
	} else if connectionType == UNIDIRECTIONAL {
		initializationPacket = string(UNIDIRECTIONAL)
	} else {
		fmt.Printf("error: travelWay protocol requires either \"bidirectional ip:port\" or \"unidirectional\", got %s\n", connectionType)
		return
	}

	initializationPacket += " " + string(travelWays.city.name)

	// send first packet
	_, err := connection.Write([]byte(initializationPacket))
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

	// make sure we aren't already connected to the city
	if _, alreadyExist := travelWays.city.outboundTravelWays[remoteCityName]; alreadyExist {
		fmt.Printf("error: travelWay from %s to %s already exists\n", travelWays.city.name, remoteCityName)
		return
	}

	// add travelWay to city
	channel := make(chan *Merchant, 100)
	travelWays.city.outboundTravelWays[remoteCityName] = channel

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
