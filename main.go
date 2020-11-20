/*
This middleware made for auth system that randomly generate access tokens, which used later for accessing secure content. Since there is no pre-defined token value, naive approach without middleware (or if middleware use only request payloads) will fail, because replayed server have own tokens, not synced with origin. To fix this, our middleware should take in account responses of replayed and origin server, store `originalToken -> replayedToken` aliases and rewrite all requests using this token to use replayed alias. See `middleware_test.go#TestTokenMiddleware` test for examples of using this middleware.
How middleware works:
                   Original request      +--------------+
+-------------+----------STDIN---------->+              |
|  Gor input  |                          |  Middleware  |
+-------------+----------STDIN---------->+              |
                   Original response     +------+---+---+
                                                |   ^
+-------------+    Modified request             v   |
| Gor output  +<---------STDOUT-----------------+   |
+-----+-------+                                     |
      |                                             |
      |            Replayed response                |
      +------------------STDIN----------------->----+
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	if len(os.Args) != 2 {
		panic("请输入日志服务器链接")
	}
	for scanner.Scan() {
		encoded := scanner.Bytes()
		buf := make([]byte, len(encoded)/2)
		hex.Decode(buf, encoded)

		process(buf)
	}
}

func process(buf []byte) {
	// First byte indicate payload type, possible values:
	//  1 - Request
	//  2 - Response
	//  3 - ReplayedResponse
	payloadType := buf[0]
	headerSize := bytes.IndexByte(buf, '\n') + 1
	header := buf[:headerSize-1]
	//Debug("----------header--------\n", string(header),
	//	"\n--------header-end----------\n")
	// Header contains space separated values of: request type, request id, and request start time (or round-trip time for responses)
	meta := bytes.Split(header, []byte(" "))

	// For each request you should receive 3 payloads (request, response, replayed response) with same request id
	reqID := string(meta[1])
	time := string(meta[2])
	payload := buf[headerSize:]
	payloadArr := bytes.Split(payload, []byte("\r\n\r\n"))
	//if len(payloadArr) > 1 {
	//	Debug("----------payloadHeaderd--------\n",
	//		string(payloadArr[0]),
	//		"\n--------payloadHeaderd-end----------\n")
	//	Debug("----------payloadBody--------\n",
	//		string(payloadArr[1]),
	//		"\n--------payloadBody-end----------\n")
	//}
	// 消噪
	if string(payload) == "" {
		return
	}
	Debug("payload", string(buf))
	switch payloadType {
	case '1': // Request
		hostSize := bytes.Index(payload, []byte("\r\n"))

		m := map[string]interface{}{
			"time":  time,
			"type":  string(payloadType),
			"reqID": reqID,
			"all":   string(buf),
			"path":  string(payload[:hostSize]),
		}
		if len(payloadArr) > 1 && string(payloadArr[1]) != "" {
			m["body"] = string(payloadArr[1])
		}

		write(m)
		// Emitting data back
		os.Stdout.Write(encode(buf))
	case '2': // Original response
		m := map[string]interface{}{
			"type":  string(payloadType),
			"reqID": reqID,
			"all":   string(buf),
		}
		if len(payloadArr) > 1 && string(payloadArr[1]) != "" {
			m["body"] = string(payloadArr[1])
		}
		write(m)
	case '3': // Replayed response
		m := map[string]interface{}{
			"type":  string(payloadType),
			"reqID": reqID,
			"all":   string(buf),
		}
		if len(payloadArr) > 1 && string(payloadArr[1]) != "" {
			m["body"] = string(payloadArr[1])
		}
		write(m)
	}
}

func encode(buf []byte) []byte {
	dst := make([]byte, len(buf)*2+1)
	hex.Encode(dst, buf)
	dst[len(dst)-1] = '\n'

	return dst
}

func Debug(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}

func write(m map[string]interface{}) {
	j, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)
	}
	_, _ = http.Post(os.Args[1]+"/write", "application/json",
		bytes.NewReader(j))
}
