package xchainge

import (
	"xchainge/types"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

// AminoCdc amino编码类
var AminoCdc = amino.NewCodec()

func init() {
	AminoCdc.RegisterInterface((*types.IPayload)(nil), nil)
	AminoCdc.RegisterConcrete(&types.Bottle{}, "bottle", nil)
	AminoCdc.RegisterConcrete(&types.Message{}, "message", nil)
	AminoCdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PubKey{}, "ed25519/pubkey", nil)
	AminoCdc.RegisterInterface((*crypto.PrivKey)(nil), nil)
	AminoCdc.RegisterConcrete(ed25519.PrivKey{}, "ed25519/privkey", nil)
}
