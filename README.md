# zenoh-go

[![Go Reference](https://pkg.go.dev/badge/github.com/wind-c/zenoh-go.svg)](https://pkg.go.dev/github.com/wind-c/zenoh-go)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)

Go bindings for zenoh-c: zero-copy, pub/sub, queryable protocol supporting router, peer, and client modes.

## Why zenoh-go?

zenoh-go provides native Go bindings for [zenoh-c](https://github.com/eclipse-zenoh/zenoh-c), offering a clean and idiomatic Go API for the zenoh protocol. Here's why you should consider zenoh-go:

- **Zero-Copy Performance**: zenoh's pub/sub and queryable patterns minimize data copying, perfect for high-performance edge computing
- **Unified Communication**: Combines pub/sub, query/response, and storage in a single protocol
- **Edge-Ready**: Optimized for constrained environments with efficient resource usage
- **Cross-Platform**: Runs on Linux, macOS, and Windows, with support for ARM64 (including Raspberry Pi, RK3576)
- **Distributed Discovery**: Automatic peer discovery via multicast, no configuration required
- **Shared Memory**: Efficient inter-process communication via shared memory protocol

## Features

### Core Features
- **Publish/Subscribe**: Pub/sub with configurable reliability and congestion control
- **Query/Queryable**: Request-response pattern for client-server interactions
- **Key Expressions**: Wildcard-based topic matching with set operations
- **Encoding Support**: Built-in support for text, JSON, binary, and custom encodings

### Advanced Features
- **Session Management**: Client and Peer modes
- **Transport**: UDP multicast, TCP, QUIC
- **Scout/Discovery**: Automatic peer and router discovery
- **Matching Status**: Track subscriber/publisher matching state
- **Shared Memory**: Zero-copy SHM protocol (requires zenoh-c with Z_FEATURE_SHM)

### Memory Management
- Explicit ownership model matching zenoh-c
- All owned types require explicit `Drop()` calls
- Clear distinction between owned and loaned types

## Architecture

```
┌────────────────────────────────────────────────────────────┐
│                        User Application                    │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────┐
│                     pkg/zenoh (Public API)                 │
│  ┌─────────┐ ┌──────────┐ ┌───────────┐ ┌──────────────┐   │
│  │ Config  │ │ Session  │ │ Publisher │ │ Subscriber   │   │
│  └─────────┘ └──────────┘ └───────────┘ └──────────────┘   │
│  ┌─────────┐ ┌──────────┐ ┌───────────┐ ┌──────────────┐   │
│  │ KeyExpr │ │ Encoding │ │  Query    │ │ Queryable    │   │
│  └─────────┘ └──────────┘ └───────────┘ └──────────────┘   │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────┐
│                  internal/cgo (CGO Bindings)               │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              zenoh_c.go (zenoh-c wrapper)            │  │
│  └──────────────────────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────┐
│                      zenoh-c (C Library)                   │
└────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌────────────────────────────────────────────────────────────┐
│                    Network (UDP/TCP/QUIC/SHM)              │
└────────────────────────────────────────────────────────────┘
```

## Project Structure

```
zenoh-go/
├── cmd/                         # Example applications
│   └── examples/
│       ├── pub/                 # Publisher example
│       ├── sub/                 # Subscriber example
│       ├── query/               # Query example (client)
│       ├── queryable/           # Queryable example (server)
│       ├── quic-pub/            # QUIC publisher example
│       ├── quic-sub/            # QUIC subscriber example
│       └── shm-pub/             # Shared memory example
├── pkg/zenoh/                   # Public Go API
│   ├── config.go                # Configuration management
│   ├── session.go               # Session handling
│   ├── publisher.go             # Publisher API
│   ├── subscriber.go            # Subscriber API
│   ├── query.go                 # Query API
│   ├── queryable.go             # Queryable API
│   ├── keyexpr.go               # Key expression handling
│   ├── encoding.go              # Encoding definitions
│   ├── bytes.go                 # Bytes serialization
│   ├── shm.go                   # Shared memory
│   ├── scout.go                 # Discovery
│   └── types.go                 # Core type definitions
├── internal/
│   └── cgo/
│       ├── zenoh_c.go           # CGO bindings
│       └── include/             # Header files
├── lib/                         # Platform-specific libraries
│   ├── windows/
│   │   ├── libzenohc.a          # Static library
│   │   └── zenohc.dll           # Dynamic library
│   ├── linux/
│   │   ├── libzenohc.a          # Static library
│   │   └── libzenohc.so         # Shared library
│   └── darwin/
│       ├── libzenohc.a          # Static library
│       └── libzenohc.dylib      # Shared library
├── bin/                         # Compiled binaries
├── go.mod                       # Go module definition
├── README.md                    # This file
└── LICENSE                      # Apache 2.0 License
```

## Installation

### Prerequisites

- Go 1.21 or later
- C compiler (GCC, Clang, or MSVC)
- zenoh-c 1.7.x library

### Build

```powershell
# Build all examples
go build -o bin/ ./cmd/examples/...

# Or build specific example
go build -o bin/zenoh-pub.exe ./cmd/examples/pub/
```

### Run

```powershell
# Terminal 1
./bin/zenoh-sub.exe

# Terminal 2
./bin/zenoh-pub.exe
```

### Linux / macOS

```bash
# Build
go build -o bin/zenoh-pub ./cmd/examples/pub/

# Run
./bin/zenoh-sub   # Terminal 1
./bin/zenoh-pub   # Terminal 2
```

### Cross-Compilation

> Pre-built zenoh-c libraries available at [zenoh-c releases](https://github.com/eclipse-zenoh/zenoh-c/releases)

```bash
# ARM64 example
GOOS=linux GOARCH=arm64 CGO_ENABLED=1 \
  CC=aarch64-linux-gnu-gcc \
  go build -o bin/zenoh-pub-arm64 ./cmd/examples/pub/
```

## Using zenoh-go as a Library

### Add Dependency

```bash
go get github.com/wind-c/zenoh-go/pkg/zenoh
```

### Import in Your Code

```go
import "github.com/wind-c/zenoh-go/pkg/zenoh"
```

### Build Requirements

| Requirement | Description |
|------------|-------------|
| CGO Enabled | `CGO_ENABLED=1` (enabled by default) |
| zenoh-c Library | Required at compile time |

### Library Location

The build will automatically link against the zenoh-c library in:

```
lib/
├── windows/   # For Windows builds
├── linux/     # For Linux builds
└── darwin/   # For macOS builds
```

> **Note**: Pre-built zenoh-c libraries are available at [zenoh-c releases](https://github.com/eclipse-zenoh/zenoh-c/releases)

### Build Your Application
For Windows Platform, other platforms similarly:
```powershell
# Build (library is automatically linked from lib/ directory)
go build -o myapp.exe .

# Or specify library path manually
$env:CGO_LDFLAGS = "-L C:\path\to\zenohc\lib"
go build -o myapp.exe .
```

### Runtime

If using dynamic library (zenohc.dll):
- Place the DLL in the same directory as your executable, OR
- Add the DLL directory to PATH environment variable

If using static library (libzenohc.a): No additional runtime configuration needed.

## Quick Start

### Publisher Example

```go
package main

import (
    "log"
    "time"

    "github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
    // Create default configuration (peer mode)
    config, err := zenoh.NewDefaultConfig()
    if err != nil {
        log.Fatal("Failed to create config: ", err)
    }
    defer config.Drop()

    // Configure peer mode
    config.InsertJSON5("mode", "\"peer\"")
    config.InsertJSON5("transport/unicast/enabled", "true")
    config.InsertJSON5("transport/multicast/enabled", "true")

    // Open session
    session, err := zenoh.Open(config)
    if err != nil {
        log.Fatal("Failed to open session: ", err)
    }
    defer session.Drop()

    // Declare publisher
    keyExpr := "demo/example/test"
    publisher, err := zenoh.DeclarePublisherWithKeyExpr(session, keyExpr)
    if err != nil {
        log.Fatal("Failed to declare publisher: ", err)
    }
    defer publisher.Undeclare()

    log.Printf("Publisher declared on: %s", keyExpr)

    // Publish messages
    for i := 0; i < 10; i++ {
        value := []byte("Hello from zenoh-go!")
        err := publisher.Put(value, zenoh.TextPlain())
        if err != nil {
            log.Printf("Failed to publish: %v", err)
        } else {
            log.Printf("Published: %s", value)
        }
        time.Sleep(1 * time.Second)
    }
}
```

### Subscriber Example

```go
package main

import (
    "log"
    "os"
    "os/signal"

    "github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
    // Create default configuration
    config, err := zenoh.NewDefaultConfig()
    if err != nil {
        log.Fatal("Failed to create config: ", err)
    }
    defer config.Drop()

    // Configure peer mode
    config.InsertJSON5("mode", "\"peer\"")
    config.InsertJSON5("transport/unicast/enabled", "true")
    config.InsertJSON5("transport/multicast/enabled", "true")

    // Open session
    session, err := zenoh.Open(config)
    if err != nil {
        log.Fatal("Failed to open session: ", err)
    }
    defer session.Drop()

    // Declare subscriber
    keyExpr := "demo/example/test"
    subscriber, err := zenoh.DeclareSubscriber(session, keyExpr, func(sample zenoh.Sample) {
        log.Printf("Received - Key: %s, Value: %s", sample.KeyExpr, string(sample.Payload))
    })
    if err != nil {
        log.Fatal("Failed to declare subscriber: ", err)
    }
    defer subscriber.Undeclare()

    log.Printf("Subscriber declared on: %s", keyExpr)
    log.Println("Waiting for messages... Press Ctrl+C to exit")

    // Wait for interrupt
    sigCh := make(chan os.Signal, 1)
    os.Signal.Notify(sigCh, os.Interrupt)
    <-sigCh
}
```

### Queryable Example (Server)

```go
package main

import (
    "log"
    "os"
    "os/signal"

    "github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
    config, err := zenoh.NewDefaultConfig()
    if err != nil {
        log.Fatal("Failed to create config: ", err)
    }
    defer config.Drop()

    config.InsertJSON5("mode", "\"peer\"")

    session, err := zenoh.Open(config)
    if err != nil {
        log.Fatal("Failed to open session: ", err)
    }
    defer session.Drop()

    // Declare queryable - responds to queries
    keyExpr := "demo/**"
    queryable, err := zenoh.DeclareQueryable(session, keyExpr, func(query zenoh.Query) {
        log.Printf("Received query - Key: %s", query.KeyExpr())
        
        // Reply with data
        value := "Hello from Queryable!"
        query.Reply(query.KeyExpr(), []byte(value), zenoh.EncodingTextPlain)
    })
    if err != nil {
        log.Fatal("Failed to declare queryable: ", err)
    }
    defer queryable.Undeclare()

    log.Printf("Queryable declared on: %s", keyExpr)
    
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt)
    <-sigCh
}
```

### Query Example (Client)

```go
package main

import (
    "log"

    "github.com/wind-c/zenoh-go/pkg/zenoh"
)

func main() {
    config, err := zenoh.NewDefaultConfig()
    if err != nil {
        log.Fatal("Failed to create config: ", err)
    }
    defer config.Drop()

    session, err := zenoh.Open(config)
    if err != nil {
        log.Fatal("Failed to open session: ", err)
    }
    defer session.Drop()

    // Query data from queryables
    selector := "demo/**"
    iterator, err := zenoh.GetWithIterator(session, selector)
    if err != nil {
        log.Fatal("Failed to get: ", err)
    }

    log.Printf("Querying: %s", selector)

    for iterator.Next() {
        reply := iterator.Reply()
        if reply.IsOk() {
            log.Printf("Reply - Key: %s, Value: %s", reply.KeyExpr(), string(reply.Value()))
        } else {
            log.Printf("Error: %s", reply.Error())
        }
    }
}
```

## Development

### Running Tests

```bash
cd zenoh-go
go test ./pkg/zenoh/...
```

### Building

```bash
# Build all examples
cd zenoh-go
go build ./cmd/examples/...

# Build specific examples
go build -o bin/zenoh-pub.exe ./cmd/examples/pub/
go build -o bin/zenoh-sub.exe ./cmd/examples/sub/
go build -o bin/query.exe ./cmd/examples/query/
go build -o bin/queryable.exe ./cmd/examples/queryable/
go build -o bin/quic-pub.exe ./cmd/examples/quic-pub/
go build -o bin/quic-sub.exe ./cmd/examples/quic-sub/
go build -o bin/shm-pub-pub.exe ./cmd/examples/shm-pub-pub/
```

### Running Query/Queryable

```powershell
# Terminal 1: Start queryable (server)
.\bin\queryable.exe

# Terminal 2: Start query (client)
.\bin\query.exe
```

### Running QUIC Transport

```powershell
# Terminal 1: Start QUIC subscriber
.\bin\quic-sub.exe

# Terminal 2: Start QUIC publisher
.\bin\quic-pub.exe
```

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

```
Copyright 2024 Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
