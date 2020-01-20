package iport

import (
	"fmt"
	"isesol.com/iport/agent"
	"testing"
)

func TestCloudLogin(t *testing.T) {
	resp := agent.cloudLogin("https://api.i5sesol.com/agentServer/verify/mqtt", "A131420035", "1", "1")
	fmt.Printf("%v", resp)
}
