package testutil

import (
	"github.com/tak1827/light-nft-indexer/data"
	"github.com/tak1827/light-nft-indexer/store"
)

const (
	TestChainEndpoint = "http://127.0.0.1:8545"
	TestFaucet        = "d1c71e71b06e248c8dbe94d49ef6d6b0d64f5d71b1e33a0f39e14dadb070304a"
	TestFaucetAddress = "0xE3b0DE0E4CA5D3CB29A9341534226C4D31C9838f"
	TestPri1          = "8179ce3d00ac1d1d1d38e4f038de00ccd0e0375517164ac5448e3acc847acb34"
	TestAddress1      = "0x26fa9f1a6568b42e29b1787c403B3628dFC0C6FE"
	TestPri2          = "df38daebd09f56398cc8fd699b72f5ea6e416878312e1692476950f427928e7d"
	TestAddress2      = "0x31a6EE302c1E7602685c86EF7a3069210Bc26670"
	TestPri3          = "97d12403ffc2faa3660730ae58bca14a894ebd78b4d8207d22083554ae96be5c"
	TestAddress3      = "0xa52ce7A3B18095800ed1f550065DF9Cd5ca5ce9f"
)

func InitDB(ds ...data.StorableData) store.DB {
	var (
		db, _ = store.NewPebbleDB("", true, nil)
		batch = db.Batch()
	)
	defer batch.Close()

	batch.Put(ds...)
	batch.Commit()

	return &db
}
