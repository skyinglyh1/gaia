module github.com/cosmos/gaia

go 1.12

replace (
	github.com/cosmos/cosmos-sdk v0.38.4 => ./../cosmos-sdk
	github.com/tendermint/tendermint v0.33.3 => ./../../tendermint/tendermint
)

require (
	github.com/coreos/go-etcd v2.0.0+incompatible // indirect
	github.com/cosmos/cosmos-sdk v0.38.4
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d // indirect
	github.com/cpuguy83/go-md2man v1.0.10 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/golang/mock v1.3.1 // indirect
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/otiai10/copy v1.0.2
	github.com/otiai10/curr v0.0.0-20190513014714-f5a3d24e5776 // indirect
	github.com/pkg/errors v0.9.1
	github.com/polynetwork/cosmos-poly-module v0.0.0-20200616094009-ee7d1ce2b07d
	github.com/rakyll/statik v0.1.6 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20190704165056-9c2d0518ed81 // indirect
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.6.1
	github.com/stumble/gorocksdb v0.0.3 // indirect
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.1
	github.com/ugorji/go/codec v0.0.0-20181204163529-d75b2dcb6bc8 // indirect
	golang.org/x/text v0.3.2 // indirect
)
