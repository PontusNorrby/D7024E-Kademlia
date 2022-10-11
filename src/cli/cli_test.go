package cli

import (
	"github.com/PontusNorrby/D7024E-Kademlia/src/kademlia"
	"testing"
)

func TestStartCLI(t *testing.T) {
	type args struct {
		nodeShutDown func()
		kademlia     *kademlia.Kademlia
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartCLI(tt.args.nodeShutDown, tt.args.kademlia)
		})
	}
}

func Test_get(t *testing.T) {
	type args struct {
		input    func() string
		kademlia *kademlia.Kademlia
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			get(tt.args.input, tt.args.kademlia)
		})
	}
}

func Test_helpList(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helpList()
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		userInput    func() string
		nodeShutDown func()
		kademlia     *kademlia.Kademlia
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			run(tt.args.userInput, tt.args.nodeShutDown, tt.args.kademlia)
		})
	}
}

func Test_store(t *testing.T) {
	type args struct {
		input func() string
		Store func(data []byte) ([]*kademlia.KademliaID, string)
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store(tt.args.input, tt.args.Store)
		})
	}
}

func Test_userInput(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := userInput(); got != tt.want {
				t.Errorf("userInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
