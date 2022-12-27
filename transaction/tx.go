package transaction

import (
	"context"
	"fmt"

	belltypes "github.com/CosmWasm/wasmd/x/bellchain/types"
	"github.com/gobuffalo/packr/v2/file/resolver/encoding/hex"

	cliTx "github.com/cosmos/cosmos-sdk/client/tx"
	secp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simapp "github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	connection "interact/connection"
)

func Transaction(ctx context.Context) error {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("wasm", "wasm1f82nxmtaw0n4ljzldmfqzz5g2dycme8t4f35u3")
	creater, err := sdk.AccAddressFromBech32("wasm1f82nxmtaw0n4ljzldmfqzz5g2dycme8t4f35u3")
	if err != nil {
		return err

	}

	encCfg := simapp.MakeTestEncodingConfig()

	txBuilder := encCfg.TxConfig.NewTxBuilder()

	chainId := "bell"
	priv := "d78dc33ba74cacc01babbaa1ad00b84efcdbb62688a4693930fa37ef2d6b5f12"
	privB, _ := hex.DecodeString(priv)
	accountSeq := uint64(2)
	accountNumber := uint64(0)
	priv1 := secp256k1.PrivKey{Key: privB}

	msg := belltypes.NewMsgKycModule(creater.String(), "bell122kb44u2u25u5u53uo25u2o2uuo", true, "bell121jbu41u4h14h41ou41o4u1ou1o4hu14ou14", "13hv1h4y1v41yh14y14g15gh1515g1g55h")
	txBuilder.SetMsgs(msg)
	txBuilder.SetGasLimit(200000)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("stake", 20)))
	txBuilder.SetMemo("add kyc info")

	//The first round: We collect all signer information. We use the "set empty signature" technique to do this
	sign := signing.SignatureV2{
		PubKey: priv1.PubKey(),
		Data: &signing.SingleSignatureData{
			SignMode:  encCfg.TxConfig.SignModeHandler().DefaultMode(),
			Signature: nil,
		},
		Sequence: accountSeq,
	}

	err = txBuilder.SetSignatures(sign)
	if err != nil {
		panic(err)
	}

	//Second round: Set all signer information, so every signer can sign.
	sign = signing.SignatureV2{}
	signerD := xauthsigning.SignerData{
		ChainID:       chainId,
		AccountNumber: accountNumber,
		Sequence:      accountSeq,
	}

	sign, err = cliTx.SignWithPrivKey(
		encCfg.TxConfig.SignModeHandler().DefaultMode(), signerD,
		txBuilder, cryptotypes.PrivKey(&priv1), encCfg.TxConfig, accountSeq,
	)
	if err != nil {
		panic(err)
	}

	err = txBuilder.SetSignatures(sign)
	if err != nil {
		panic(err)
	}

	// Generated Protobuf-encoded bytes.
	txBytes, err := encCfg.TxConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return err
	}

	// Create a connection to the gRPC server.
	grpcConn := connection.Connection()
	defer grpcConn.Close()

	// Broadcast the tx via gRPC. We create a new client for the Protobuf Tx
	// service.

	txClient := tx.NewServiceClient(grpcConn)
	//we can call the BroadcastTx method on this client
	grpcRes, err := txClient.BroadcastTx(
		ctx,
		&tx.BroadcastTxRequest{
			Mode:    tx.BroadcastMode_BROADCAST_MODE_ASYNC,
			TxBytes: txBytes,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(grpcRes.GetTxResponse())
	return nil
}
