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
	"fmt"

	"github.com/dsrvlabs/vatz/sdk"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	addr = "0.0.0.0"
	port = 9091
)

func main() {
	p := sdk.NewPlugin()
	p.Register(pluginFeature)

	ctx := context.Background()
	if err := p.Start(ctx, addr, port); err != nil {
		fmt.Println("exit")
	}
}

func pluginFeature(info, opt map[string]*structpb.Value) error {
	// TODO: Fill here.
	return nil
}
```

Then build source code.

```
~$ go mod tidy
~$ go build
```
