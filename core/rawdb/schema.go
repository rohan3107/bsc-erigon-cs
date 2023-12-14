package rawdb

import libcommon "github.com/ledgerwatch/erigon-lib/common"

var (
	genesisPrefix = []byte("ethereum-genesis-")
)

// genesisStateSpecKey = genesisPrefix + hash
func genesisStateSpecKey(hash libcommon.Hash) []byte {
	return append(genesisPrefix, hash.Bytes()...)
}
