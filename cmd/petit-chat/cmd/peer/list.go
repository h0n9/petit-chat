package peer

import (
	"bufio"
	"fmt"

	"github.com/h0n9/petit-chat/util"
)

var listCmd = util.NewCmd(
	"list",
	"show a list of peers on network",
	listFunc,
)

func listFunc(reader *bufio.Reader) error {
	peers := cli.GetPeers()
	if len(peers) == 0 {
		fmt.Printf("none\n")
		return nil
	}

	for i, peer := range peers {
		fmt.Printf("%d. %s\n", i+1, peer)
	}

	return nil
}
