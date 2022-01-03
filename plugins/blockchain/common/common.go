package common

import "github.com/meshplus/hyperbench/common"

// Invoke define need filed for invoke contract.
type Invoke struct {
	Func string        `mapstructure:"func"`
	Args []interface{} `mapstructure:"args"`
}

// Transfer define need filed for transfer.
type Transfer struct {
	From   string `mapstructure:"from"`
	To     string `mapstructure:"to"`
	Amount int64  `mapstructure:"amount"`
	Extra  string `mapstructure:"extra"`
}

// Query define need filed for query info.
type Query struct {
	Func string        `mapstructure:"func"`
	Args []interface{} `mapstructure:"args"`
}

// Option for receive options.
type Option map[string]interface{}

// Context the context in vm.
type Context string

// Statistic contains statistic time.
type Statistic struct {
	From int64 `mapstructure:"from"`
	To   int64 `mapstructure:"to"`
	TimeSeq []int64 `mapstructure:"time_seq"`
	HeightSeq []uint64 `mapstructure:"height_seq"`
	data common.Data
}
