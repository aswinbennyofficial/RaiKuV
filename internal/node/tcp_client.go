package node

import (
	"fmt"
	"io"
	"net"

	"github.com/aswinbennyofficial/raikuv/internal/storage"
	"github.com/linkedin/goavro/v2"
	"github.com/rs/zerolog"
)

type TCPServer struct {
    Address   string
    DataStore storage.Storage
    Logger    *zerolog.Logger
}

func NewTCPServer(address string, dataStore storage.Storage, logger *zerolog.Logger) *TCPServer {
    return &TCPServer{
        Address:   address,
        DataStore: dataStore,
        Logger:    logger,
    }
}

func (server *TCPServer) ListenAndServe() error {
    logger := server.Logger
    
    listener, err := net.Listen("tcp", ":"+server.Address)
    if err != nil {
        return fmt.Errorf("failed to start server: %w", err)
    }
    defer listener.Close()
    
    logger.Info().Msgf("TCP server started on port %s", server.Address)

    // Create codec once for reuse
    codec, err := server.createCodec()
    if err != nil {
        return fmt.Errorf("failed to create codec: %w", err)
    }

    for {
        conn, err := listener.Accept()
        if err != nil {
            logger.Debug().Err(err).Msg("Failed to accept connection")
            continue
        }
        go server.HandleConnection(conn, codec)
    }
}

func (server *TCPServer) createCodec() (*goavro.Codec, error) {
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
    
    return goavro.NewCodec(schema)
}

func (server *TCPServer) HandleConnection(conn net.Conn, codec *goavro.Codec) {
    defer conn.Close()
    logger := server.Logger

    buffer := make([]byte, 1024) // Increased buffer size for larger messages

    for {
        // Read data from the connection
        n, err := conn.Read(buffer)
        if err != nil {
            if err != io.EOF {
                logger.Error().Err(err).Msg("Failed to read from connection")
            }
            return
        }

        // Process the received data
        response, err := server.processRequest(buffer[:n], codec)
        if err != nil {
            logger.Error().Err(err).Msg("Failed to process request")
            return
        }

        // Send the response back to the client
        _, err = conn.Write(response)
        if err != nil {
            logger.Error().Err(err).Msg("Failed to send response")
            return
        }
    }
}

func (server *TCPServer) processRequest(data []byte, codec *goavro.Codec) ([]byte, error) {
    // Deserialize Avro data
    nativeData, _, err := codec.NativeFromBinary(data)
    if err != nil {
        return nil, fmt.Errorf("error deserializing Avro data: %w", err)
    }

    dataMap, ok := nativeData.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid Avro data format")
    }

    // Extract fields
    method := dataMap["method"].(string)
    key := dataMap["key"].(string)
    value := dataMap["value"].(string)
    errorMsg := dataMap["errorMsg"].(string)

    // Process the request based on the method
    switch method {
    case "put":
        server.DataStore.Put(key, value)
    case "get":
        storedValue, ok := server.DataStore.Get(key)
        if !ok {
            errorMsg = fmt.Sprintf("Key %s not found", key)
        } else {
            value = fmt.Sprintf("%v", storedValue)
        }
    case "pop":
        server.DataStore.Pop(key)
    default:
        errorMsg = "Unknown method"
    }

    // Prepare and serialize response
    response := map[string]interface{}{
        "method":   method,
        "key":      key,
        "value":    value,
        "errorMsg": errorMsg,
    }

    return codec.BinaryFromNative(nil, response)
}