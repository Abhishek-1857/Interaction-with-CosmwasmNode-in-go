package query

import (
	"context"
	"fmt"
	belltypes "github.com/CosmWasm/wasmd/x/bellchain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connection "interact/connection"
)	

func QueryState() error {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("wasm", "wasm1f82nxmtaw0n4ljzldmfqzz5g2dycme8t4f35u3")
	_, err := sdk.AccAddressFromBech32("wasm1f82nxmtaw0n4ljzldmfqzz5g2dycme8t4f35u3")
	if err != nil {
		return err

	}

	// Create a connection to the gRPC server.
	grpcConn := connection.Connection()
	defer grpcConn.Close()

	// This creates a gRPC client to query the x/Bellchain service.
	bellClient := belltypes.NewQueryClient(grpcConn)
	bellRes, err := bellClient.Kycs(context.Background(), &belltypes.QueryKycsRequest{})
	if err != nil {
		return err
	}

	fmt.Println(bellRes)

	return nil
}
