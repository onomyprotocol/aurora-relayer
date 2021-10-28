// This is free and unencumbered software released into the public domain.

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/cloudflare/cloudflare-go"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/jhttp"
)

type NEARSyncInfo struct {
	LatestBlockHeight int64  `json:"latest_block_height"`
	LatestBlockTime   string `json:"latest_block_time"`
}

type NEARStatus struct {
	ChainID               string       `json:"chain_id"`
	LatestProtocolVersion int          `json:"latest_protocol_version"`
	ProtocolVersion       int          `json:"protocol_version"`
	SyncInfo              NEARSyncInfo `json:"sync_info"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: sync NEAR-URL\n")
		os.Exit(64) // EX_USAGE
	}
	endpointURL := os.Args[1]

	cfAPIToken := os.Getenv("CF_API_TOKEN")
	if cfAPIToken == "" {
		fmt.Fprintf(os.Stderr, "missing %s environment variable\n", "CF_API_TOKEN")
		os.Exit(78) // EX_CONFIG
	}

	cfAccountID := os.Getenv("CF_ACCOUNT_ID")
	if cfAccountID == "" {
		fmt.Fprintf(os.Stderr, "missing %s environment variable\n", "CF_ACCOUNT_ID")
		os.Exit(78) // EX_CONFIG
	}

	cfNamespaceID := os.Getenv("CF_NAMESPACE_ID")
	if cfNamespaceID == "" {
		fmt.Fprintf(os.Stderr, "missing %s environment variable\n", "CF_NAMESPACE_ID")
		os.Exit(78) // EX_CONFIG
	}

	api, err := cloudflare.NewWithAPIToken(cfAPIToken, cloudflare.UsingAccount(cfAccountID))
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(77) // EX_NOPERM
	}

	ctx := context.Background()
	channel := jhttp.NewChannel(endpointURL, nil)
	client := jrpc2.NewClient(channel, nil)

	for {
		response, err := client.Call(ctx, "status", []string{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(69) // EX_UNAVAILABLE
		}

		var status NEARStatus
		if err := response.UnmarshalResult(&status); err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(76) // EX_PROTOCOL
		}
		blockNumber := status.SyncInfo.LatestBlockHeight
		//fmt.Println(blockNumber)

		_, err = api.WriteWorkersKV(ctx, cfNamespaceID, "eth_blockNumber", []byte(strconv.FormatInt(blockNumber, 10)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(69) // EX_UNAVAILABLE
		}
	}
}
