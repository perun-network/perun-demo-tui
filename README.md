# Perun Demo TUI
This is a demo to showcase Perun Payment Channels for different backends.

## Usage
1. Implement the `client.DemoClient` interface (very similar to our usual payment-channel demo).
2. Initialize the demo clients and run the demo:

```go
package main

import (
	"perun.network/go-perun/client/test"
	"perun.network/perun-demo-tui/view"
)

func main() {
	// Setup Backend
	...
	
	// Initialize Demo Clients
	alice := newDemoClient("Alice", ...)
	bob := newDemoClient("Bob", ...)
    
	// Run Demo
	_ = view.RunDemo("Perun XYZ Backend Demo", []DemoClient{alice, bob})
}
```

## Controls
| Keybinds     	 | Function                         	 |
|----------------|------------------------------------|
| ` CTRL + A ` 	 | switch to left column (Party A)  	 |
| ` CTRL + B ` 	 | switch to right column (Party B) 	 |
| ` r `    	     | switch to parent page            	 |
| `q`     	      | quit                             	 |