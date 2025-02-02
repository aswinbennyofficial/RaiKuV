package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
    "github.com/linkedin/goavro/v2"
)

const (
    helpText = `
Available Commands:
  put <key>:<value>  Store a value with the given key
  get <key>          Retrieve a value by key
  pop <key>          Remove a value by key
  clear              Clear the screen
  help               Show this help message
  exit               Exit the client

Examples:
  > put mykey:hello world    # Store "hello world" with key "mykey"
  > get mykey               # Retrieve value for "mykey"
  > pop mykey              # Remove value for "mykey"
`
)

type Client struct {
    conn   net.Conn
    codec  *goavro.Codec
    reader *bufio.Reader
}

func NewClient(address string) (*Client, error) {
    // Connect to the TCP server
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return nil, fmt.Errorf("failed to connect: %w", err)
    }

    // Avro schema definition
    schema := `{
        "type": "record",
        "name": "Request",
        "fields": [
            {"name": "method", "type": "string"},
            {"name": "key", "type": "string"},
            {"name": "value", "type": "string"},
            {"name": "errorMsg", "type": "string", "default": ""}
        ]
    }`

    codec, err := goavro.NewCodec(schema)
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to create codec: %w", err)
    }

    return &Client{
        conn:   conn,
        codec:  codec,
        reader: bufio.NewReader(os.Stdin),
    }, nil
}

func (c *Client) Close() error {
    return c.conn.Close()
}

func (c *Client) Run() error {
    fmt.Println("Connected to TCP server. Type 'help' for available commands.")

    for {
        fmt.Print("\033[34m>\033[0m ") // Blue prompt
        input, err := c.reader.ReadString('\n')
        if err != nil {
            return fmt.Errorf("failed to read input: %w", err)
        }

        input = strings.TrimSpace(input)
        if input == "" {
            continue
        }

        if err := c.handleCommand(input); err != nil {
            if err.Error() == "exit" {
                return nil
            }
            fmt.Printf("\033[31mError: %v\033[0m\n", err) // Red error message
        }
    }
}

func (c *Client) handleCommand(input string) error {
    switch {
    case input == "exit":
        fmt.Println("Goodbye!")
        return fmt.Errorf("exit")
    
    case input == "help":
        fmt.Println(helpText)
        return nil
    
    case input == "clear":
        fmt.Print("\033[H\033[2J") // Clear screen
        return nil
    
    case strings.HasPrefix(input, "put "):
        return c.handlePut(input[4:])
    
    case strings.HasPrefix(input, "get "):
        return c.handleGet(input[4:])
    
    case strings.HasPrefix(input, "pop "):
        return c.handlePop(input[4:])
    
    default:
        return fmt.Errorf("unknown command. Type 'help' for available commands")
    }
}

func (c *Client) sendRequest(method, key, value string) (map[string]interface{}, error) {
    // Prepare Avro record
    record := map[string]interface{}{
        "method": method,
        "key":    key,
        "value":  value,
    }

    // Serialize and send
    avroData, err := c.codec.BinaryFromNative(nil, record)
    if err != nil {
        return nil, fmt.Errorf("serialization failed: %w", err)
    }

    if _, err := c.conn.Write(avroData); err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }

    // Read response
    response := make([]byte, 1024)
    n, err := c.conn.Read(response)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    // Deserialize response
    nativeResponse, _, err := c.codec.NativeFromBinary(response[:n])
    if err != nil {
        return nil, fmt.Errorf("failed to deserialize response: %w", err)
    }

    responseMap, ok := nativeResponse.(map[string]interface{})
    if !ok {
        return nil, fmt.Errorf("invalid response format")
    }

    return responseMap, nil
}

func (c *Client) handlePut(input string) error {
    parts := strings.SplitN(input, ":", 2)
    if len(parts) != 2 || parts[0] == "" {
        return fmt.Errorf("invalid format. Use: put key:value")
    }

    response, err := c.sendRequest("put", parts[0], parts[1])
    if err != nil {
        return err
    }

    if errMsg := response["errorMsg"].(string); errMsg != "" {
        return fmt.Errorf(errMsg)
    }

    fmt.Printf("\033[32mSuccessfully stored value for key '%s'\033[0m\n", parts[0])
    return nil
}

func (c *Client) handleGet(key string) error {
    if key == "" {
        return fmt.Errorf("key cannot be empty")
    }

    response, err := c.sendRequest("get", key, "")
    if err != nil {
        return err
    }

    if errMsg := response["errorMsg"].(string); errMsg != "" {
        return fmt.Errorf(errMsg)
    }

    fmt.Printf("\033[32mValue: %s\033[0m\n", response["value"])
    return nil
}

func (c *Client) handlePop(key string) error {
    if key == "" {
        return fmt.Errorf("key cannot be empty")
    }

    response, err := c.sendRequest("pop", key, "")
    if err != nil {
        return err
    }

    if errMsg := response["errorMsg"].(string); errMsg != "" {
        return fmt.Errorf(errMsg)
    }

    fmt.Printf("\033[32mSuccessfully removed key '%s'\033[0m\n", key)
    return nil
}

func main() {
    client, err := NewClient("localhost:3232")
    if err != nil {
        fmt.Printf("\033[31mFailed to start client: %v\033[0m\n", err)
        os.Exit(1)
    }
    defer client.Close()

    if err := client.Run(); err != nil {
        fmt.Printf("\033[31mError: %v\033[0m\n", err)
        os.Exit(1)
    }
}