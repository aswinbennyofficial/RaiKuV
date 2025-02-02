package node

import (
	"fmt"
	"net"

	"github.com/aswinbennyofficial/raikuv/internal/storage"
	"github.com/linkedin/goavro/v2"
	"github.com/rs/zerolog"
)

type TCPServer struct{
	Address string
	DataStore *storage.DataStore
	Logger *zerolog.Logger
}

func NewTCPServer(address string, dataStore *storage.DataStore, logger *zerolog.Logger)(*TCPServer){
	return &TCPServer {
		Address : address,
		DataStore: dataStore,
		Logger: logger,
	}
}

func (server *TCPServer) ListenAndServe() error{
	logger := server.Logger

	// Start listening on the address
	listener, err:= net.Listen("tcp",":"+server.Address)
	if err!=nil{
		return err
	}

	defer listener.Close()

	logger.Info().Msgf("TCP server started in port %s",server.Address)
	

	for {
		conn,err:=listener.Accept()
		if err!=nil{
			logger.Debug().Msgf("failed to accept connection: %v", err)
			continue
		}

		go server.HandleConnection(conn)
	}

}

func (server *TCPServer)HandleConnection(conn net.Conn){
	defer conn.Close()
	
	logger :=server.Logger

	schema := `{
        "type": "record",
        "name": "Request",
        "fields": [
            {"name": "method", "type": "string"},
            {"name": "key", "type": "string"},
            {"name": "value", "type": "string"},
			{"name": "errorMsg", "type": "string"}
        ]
    }`

	for {

		// Parse the schema and create codec
		codec, err := goavro.NewCodec(schema)
		if err != nil {
			logger.Println("Error creating Avro codec:", err)
			return
		}

		// Read data from the connection
		buffer := make([]byte, 128)
		n, err := conn.Read(buffer)
		if err != nil {
			logger.Println("Failed to read from connection:", err)
			return
		}

		// Deserialize Avro data
		binaryData := buffer[:n]
		nativeData,_,err := codec.NativeFromBinary(binaryData)
		if err != nil {
			logger.Println("Error deserializing Avro data:", err)
			return
		}

		// Cast to the expected map structure
		dataMap, ok := nativeData.(map[string]interface{})
		if !ok {
			logger.Println("Invalid Avro data format")
			return
		}

		// Extract method, key, and value from the deserialized data
		method := dataMap["method"].(string)
		key := dataMap["key"].(string)
		value := dataMap["value"].(string)
		errorMsg:= dataMap["errorMsg"].(string)

		// Process the request based on the method
		if method == "put" {
			// Put operation
			server.DataStore.Put(key, value)
			
		} else if method == "get" {
			// Get operation
			storedValue, ok := server.DataStore.Get(key)
			if !ok {
				logger.Error().Msgf("Key %s not found ",key)
				errorMsg = fmt.Sprintf("Key %s not found ",key) 
			}else{
				value = fmt.Sprintf("%v",storedValue)
			}
		} else if method == "pop" {
			server.DataStore.Pop(key)
		}

		// Prepare a response
		response := map[string]interface{}{
			"method": method,
			"key":    key,
			"value":  value,
			"errorMsg":  errorMsg,
		}

		// Serialize response
		avroResponse, err := codec.BinaryFromNative(nil, response)
		if err != nil {
			logger.Println("Error serializing response:", err)
			return
		}

		// Send the response back to the client
		_, err = conn.Write(avroResponse)
		if err != nil {
			logger.Println("Failed to send response:", err)
		}

	}

}

 