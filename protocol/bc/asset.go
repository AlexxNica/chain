package bc

import (
	"database/sql/driver"
	"io"

	"chain/crypto/sha3pool"
	"chain/encoding/blockchain"
)

// AssetID is the Hash256 of the issuance script for the asset and the
// initial block of the chain where it appears.
type AssetID [32]byte

func (a AssetID) String() string                { return Hash(a).String() }
func (a AssetID) MarshalText() ([]byte, error)  { return Hash(a).MarshalText() }
func (a *AssetID) UnmarshalText(b []byte) error { return (*Hash)(a).UnmarshalText(b) }
func (a *AssetID) UnmarshalJSON(b []byte) error { return (*Hash)(a).UnmarshalJSON(b) }
func (a AssetID) Value() (driver.Value, error)  { return Hash(a).Value() }
func (a *AssetID) Scan(b interface{}) error     { return (*Hash)(a).Scan(b) }

type AssetDefinition struct {
	InitialBlockID  Hash
	IssuanceProgram Program
	Data            Hash
}

func (ad *AssetDefinition) ComputeAssetID() (assetID AssetID) {
	h := sha3pool.Get256()
	defer sha3pool.Put256(h)
	writeForHash(h, *ad) // error is impossible
	h.Read(assetID[:])
	return assetID
}

func ComputeAssetID(prog []byte, initialBlockID Hash, vmVersion uint64, data Hash) AssetID {
	def := &AssetDefinition{
		InitialBlockID: initialBlockID,
		IssuanceProgram: Program{
			VMVersion: vmVersion,
			Code:      prog,
		},
		Data: data,
	}
	return def.ComputeAssetID()
}
