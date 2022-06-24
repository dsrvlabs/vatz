# Vatz SDK for plugin

**For now, `vatz` repository is private so many part of starting new project should be done by manually.**

### Step 1: Create new project

First, create new directory to store.

```
~$ mkdir <project-directory>
~$ cd <project-directory>
```

Initialize go.mod file.

```
~$ go mod init <project name>
```

### Step 2: Set GOPRIVATE

For now, `vatz` project is in private so extra configurations are required to use SDK module.

```
~$ export GOPRIVATE=github.com/dsrvlabs/vatz,github.com/dsrvlabs/vatz-proto
~$ export GIT_TERMINAL_PROMPT=1
~$ go get github.com/dsrvlabs/vatz
```

### Step 3: Start main

Create `main.go` file with below contents.

```
~$ touch main.go
```

```
package main

import (
        "flag"
        "fmt"

        pluginpb "github.com/dsrvlabs/vatz-proto/plugin/v1"
        "github.com/dsrvlabs/vatz/sdk"
        "golang.org/x/net/context"
        "google.golang.org/protobuf/types/known/structpb"
)

const (
        // Default values.
        defaultAddr = "127.0.0.1"
        defaultPort = 9091

        pluginName = "YOUR_PLUGIN_NAME"
)

var (
        addr string
        port int
)

func init() {
        flag.StringVar(&addr, "addr", defaultAddr, "IP Address(e.g. 0.0.0.0, 127.0.0.1)")
        flag.IntVar(&port, "port", defaultPort, "Port number, defulat 9091")

        flag.Parse()
}

func main() {
        p := sdk.NewPlugin(pluginName)
        p.Register(pluginFeature)

        ctx := context.Background()
        if err := p.Start(ctx, addr, port); err != nil {
                fmt.Println("exit")
        }
}

func pluginFeature(info, option map[string]*structpb.Value) (sdk.CallResponse, error) {
        // TODO: Fill here.
        ret := sdk.CallResponse{
                FuncName:   "YOUR_FUNCTION_NAME",
                Message:    "YOUR_MESSAGE_CONTENTS",
                Severity:   pluginpb.SEVERITY_UNKNOWN,
                State:      pluginpb.STATE_NONE,
                AlertTypes: []pluginpb.ALERT_TYPE{pluginpb.ALERT_TYPE_DISCORD},
        }

        return ret, nil
}
```

Then build source code.

```
~$ go mod tidy
~$ go build
```
