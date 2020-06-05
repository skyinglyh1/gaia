package types

import (
	"fmt"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryContainToContract struct {
	KeyStore       string
	ToContractAddr []byte
	FromChainId    uint64
}

func NewQueryContainToContract(keystore string, toContractAddr []byte, fromChainId uint64) QueryContainToContract {
	return QueryContainToContract{keystore, toContractAddr, fromChainId}
}

type QueryContainToContractRes struct {
	KeyStore string
	Exist    bool
	Info     string
}

func (this QueryContainToContractRes) String() string {
	return fmt.Sprintf(`
  KeyStore:				%s,
  Exist:				%t,
  Info:					%s,
`, this.KeyStore, this.Exist, this.Info)
}
