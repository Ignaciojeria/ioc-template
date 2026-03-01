# main

> Application entry point - loads IoC dependencies

## cmd/api/main.go

```go
package main

import (
	"log"
	"os"
	"strings"

	"archetype"

	"github.com/Ignaciojeria/ioc"
)

func main() {
	os.Setenv("VERSION", strings.TrimSpace(archetype.Version))

	if err := ioc.LoadDependencies(); err != nil {
		log.Fatal("Failed to load dependencies:", err)
	}
}
```
