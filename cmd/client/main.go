package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"github.com/linkedin/goavro/v2"
)

func main() {
	// Connect to the TCP server
	conn, err := net.Dial("tcp", "localhost:3232")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	// Avro schema definition for the message structure
	schema := `{
		"type": "record",
		"name": "Request",
		"fields": [
			{"name": "method", "type": "string"},
			{"name": "key", "type": "string"},
			{"name": "value", "type": "string"},
			{"name": "errorMsg", "type": "string", "default":""}
		]
	}`

	// Create Avro codec from schema
	codec, err := goavro.NewCodec(schema)
	if err != nil {
		fmt.Println("Error creating Avro codec:", err)
		return
	}

	fmt.Println("Connected to TCP server. Type `put key:value` to store data or `get key` to retrieve.")

	reader := bufio.NewReader(os.Stdin)

	for {
		// Read input from user
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Exit condition
		if input == "exit" {
			fmt.Println("Exiting client...")
			break
		}

		// Prepare data for Avro serialization
		var method, key, value string
		parts := strings.SplitN(input, " ", 2)

		if len(parts) == 2 {
			method = parts[0]
			kvParts := strings.SplitN(parts[1], ":", 2)
			key = kvParts[0]
			if len(kvParts) > 1 {
				value = kvParts[1]
			}
		} else {
			method = parts[0]
			key = parts[1]
		}

		// Prepare Avro record for the request
		record := map[string]interface{}{
			"method": method,
			"key":    key,
			"value":  value,
			
		}

		// Serialize the request into Avro binary format
		avroData, err := codec.BinaryFromNative(nil, record)
		if err != nil {
			fmt.Println("Error serializing data:", err)
			return
		}

		// Send the serialized data to the server
		_, err = conn.Write(avroData)
		if err != nil {
			fmt.Println("Error writing to server:", err)
			return
		}

		// Read the server's response
		response := make([]byte, 1024)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		// Deserialize the server's Avro response
		nativeResponse,_, err := codec.NativeFromBinary(response[:n])
		if err != nil {
			fmt.Println("Error deserializing response:", err)
			return
		}

		// Print the response
		responseMap, ok := nativeResponse.(map[string]interface{})
		if !ok {
			fmt.Println("Error: unable to assert response to map")
			continue
		}

		if responseMap["errorMsg"]!=""{
			fmt.Println(responseMap["errorMsg"])
		}else{
			fmt.Printf("%v\n",responseMap)
		}
	}
}
