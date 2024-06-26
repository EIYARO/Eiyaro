package consensus

import (
	"encoding/binary"
	"math"
	"strings"

	"eiyaro/protocol/bc"
)

// consensus variables
const (
	// MaxBlockGas Max gas that one block contains
	MaxBlockGas      = uint64(10000000)
	VMGasRate        = int64(200)
	StorageGasRate   = int64(1)
	MaxGasAmount     = int64(200000)
	DefaultGasCredit = int64(30000)

	// CoinbasePendingBlockNumber config parameter for coinbase reward
	CoinbasePendingBlockNumber = uint64(50)
	subsidyReductionInterval   = uint64(175200)
	baseSubsidy                = uint64(100000000000)
	InitialBlockSubsidy        = uint64(21000000000000000)
	subsidyReductionRate       = 0.1

	// BlocksPerRetarget config for pow mining
	BlocksPerRetarget     = uint64(1000)
	TargetSecondsPerBlock = uint64(180)
	SeedPerRetarget       = uint64(256)

	// MaxTimeOffsetSeconds is the maximum number of seconds a block time is allowed to be ahead of the current time
	MaxTimeOffsetSeconds = uint64(10 * 60)
	MedianTimeBlocks     = 8

	PayToWitnessPubKeyHashDataSize = 20
	PayToWitnessScriptHashDataSize = 32
	CoinbaseArbitrarySizeLimit     = 128

	EYAlias = "EY"
)

// EYAssetID is EY's asset id, the soul asset of Eiyaro
var EYAssetID = &bc.AssetID{
	V0: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V1: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V2: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
	V3: binary.BigEndian.Uint64([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}),
}

// InitialSeed is SHA3-256 of Byte[0^32]
var InitialSeed = &bc.Hash{
	V0: uint64(11412844483649490393),
	V1: uint64(4614157290180302959),
	V2: uint64(1780246333311066183),
	V3: uint64(9357197556716379726),
}

// EYDefinitionMap is the ....
var EYDefinitionMap = map[string]interface{}{
	"name":        EYAlias,
	"symbol":      EYAlias,
	"decimals":    8,
	"description": `Eiyaro Official Issue`,
}

// BlockSubsidy calculates the coinbase reward for a given block height, considering a 10% reduction every subsidyReductionInterval blocks.
func BlockSubsidy(height uint64) uint64 {
	if height == 0 {
		return InitialBlockSubsidy
	}

	// Determine the number of full subsidy reduction intervals that have passed
	intervalsPassed := height / subsidyReductionInterval

	// Calculate the subsidy reduction factor based on the number of intervals passed
	reductionFactor := math.Pow(1-subsidyReductionRate, float64(intervalsPassed))

	// Apply the reduction factor to the base subsidy
	subsidy := uint64(math.Round(float64(baseSubsidy) * reductionFactor))

	return subsidy
}

// IsBech32SegwitPrefix returns whether the prefix is a known prefix for segwit
// addresses on any default or registered network.  This is used when decoding
// an address string into a specific address type.
func IsBech32SegwitPrefix(prefix string, params *Params) bool {
	prefix = strings.ToLower(prefix)
	return prefix == params.Bech32HRPSegwit+"1"
}

// Checkpoint identifies a known good point in the block chain.  Using
// checkpoints allows a few optimizations for old blocks during initial download
// and also prevents forks from old blocks.
type Checkpoint struct {
	Height uint64
	Hash   bc.Hash
}

// Params store the config for different network
type Params struct {
	// Name defines a human-readable identifier for the network.
	Name            string
	Bech32HRPSegwit string
	// DefaultPort defines the default peer-to-peer port for the network.
	DefaultPort string

	// DNSSeeds defines a list of DNS seeds for the network that are used
	// as one method to discover peers.
	DNSSeeds    []string
	Checkpoints []Checkpoint
}

// ActiveNetParams is ...
var ActiveNetParams = MainNetParams

// NetParams is the correspondence between chain_id and Params
var NetParams = map[string]Params{
	"mainnet": MainNetParams,
	"wisdom":  TestNetParams,
	"solonet": SoloNetParams,
}

// MainNetParams is the config for production
var MainNetParams = Params{
	Name:            "main",
	Bech32HRPSegwit: "ey",
	DefaultPort:     "46657",
	DNSSeeds:        []string{"mainnetseed.eiyaro.org"},
	Checkpoints: []Checkpoint{
		{10000, bc.NewHash([32]byte{0x93, 0xe1, 0xeb, 0x78, 0x21, 0xd2, 0xb4, 0xad, 0x0f, 0x5b, 0x1c, 0xea, 0x82, 0xe8, 0x43, 0xad, 0x8c, 0x09, 0x9a, 0xb6, 0x5d, 0x8f, 0x70, 0xc5, 0x84, 0xca, 0xa2, 0xdd, 0xf1, 0x74, 0x65, 0x2c})},
		{20000, bc.NewHash([32]byte{0x7d, 0x38, 0x61, 0xf3, 0x2c, 0xc0, 0x03, 0x81, 0xbb, 0xcd, 0x9a, 0x37, 0x6f, 0x10, 0x5d, 0xfe, 0x6f, 0xfe, 0x2d, 0xa5, 0xea, 0x88, 0xa5, 0xe3, 0x42, 0xed, 0xa1, 0x17, 0x9b, 0xa8, 0x0b, 0x7c})},
		{30000, bc.NewHash([32]byte{0x32, 0x36, 0x06, 0xd4, 0x27, 0x2e, 0x35, 0x24, 0x46, 0x26, 0x7b, 0xe0, 0xfa, 0x48, 0x10, 0xa4, 0x3b, 0xb2, 0x40, 0xf1, 0x09, 0x51, 0x5b, 0x22, 0x9f, 0xf3, 0xc3, 0x83, 0x28, 0xaa, 0x4a, 0x00})},
		{40000, bc.NewHash([32]byte{0x7f, 0xe2, 0xde, 0x11, 0x21, 0xf3, 0xa9, 0xa0, 0xee, 0x60, 0x8d, 0x7d, 0x4b, 0xea, 0xcc, 0x33, 0xfe, 0x41, 0x25, 0xdc, 0x2f, 0x26, 0xc2, 0xf2, 0x9c, 0x07, 0x17, 0xf9, 0xe4, 0x4f, 0x9d, 0x46})},
		{50000, bc.NewHash([32]byte{0x5e, 0xfb, 0xdf, 0xf5, 0x35, 0x38, 0xa6, 0x0b, 0x75, 0x32, 0x02, 0x61, 0x83, 0x54, 0x34, 0xff, 0x3e, 0x82, 0x2e, 0xf8, 0x64, 0xae, 0x2d, 0xc7, 0x6c, 0x9d, 0x5e, 0xbd, 0xa3, 0xd4, 0x50, 0xcf})},
		{62000, bc.NewHash([32]byte{0xd7, 0x39, 0x8f, 0x23, 0x57, 0xf9, 0x4c, 0xa0, 0x28, 0xa7, 0x00, 0x2b, 0x53, 0x9e, 0x51, 0x2d, 0x3e, 0xca, 0xc9, 0x22, 0x59, 0xfc, 0xd0, 0x3f, 0x67, 0x1a, 0x0a, 0xb1, 0x02, 0xbf, 0x2b, 0x03})},
		{72000, bc.NewHash([32]byte{0x66, 0x02, 0x31, 0x19, 0xf1, 0x60, 0x35, 0x61, 0xa4, 0xf1, 0x38, 0x04, 0xcc, 0xe4, 0x59, 0x8f, 0x55, 0x39, 0xba, 0x22, 0xf2, 0x6d, 0x90, 0xbf, 0xc1, 0x87, 0xef, 0x98, 0xcc, 0x70, 0x4d, 0x94})},
		{83700, bc.NewHash([32]byte{0x7f, 0x26, 0xc9, 0x11, 0xe8, 0x46, 0xd0, 0x6e, 0x36, 0xbb, 0xac, 0xce, 0x99, 0xa2, 0x19, 0x89, 0x3f, 0xf7, 0x84, 0x2a, 0xcb, 0x44, 0x7f, 0xbb, 0x0e, 0x3b, 0xa3, 0x68, 0xd6, 0x2b, 0xe8, 0x0d})},
		{93700, bc.NewHash([32]byte{0x70, 0x44, 0x70, 0xe5, 0xb3, 0x9b, 0xd3, 0x67, 0x19, 0x20, 0x08, 0x42, 0x1b, 0x59, 0xe8, 0xdc, 0xb5, 0xbb, 0xb9, 0x2d, 0xd3, 0xdc, 0x28, 0x4e, 0xcb, 0x7b, 0x0b, 0xbf, 0x21, 0x51, 0xe1, 0xba})},
		{106600, bc.NewHash([32]byte{0x31, 0x15, 0x2b, 0x00, 0xd4, 0x07, 0xe1, 0xa7, 0x06, 0xe1, 0xae, 0x2e, 0x98, 0x69, 0x8f, 0x47, 0xff, 0x44, 0x97, 0x01, 0xa7, 0x9e, 0x08, 0xdb, 0xeb, 0x0f, 0x1f, 0x5a, 0xdd, 0xf5, 0x26, 0xb9})},
		{116600, bc.NewHash([32]byte{0x08, 0xeb, 0xf7, 0x6c, 0x27, 0xed, 0x81, 0xe7, 0xe7, 0xfe, 0x13, 0xca, 0x80, 0x71, 0x29, 0x26, 0x28, 0x72, 0x25, 0xa5, 0x2a, 0xa0, 0x36, 0x30, 0x58, 0xaa, 0x58, 0xc6, 0xdd, 0xf2, 0xa0, 0xe7})},
		{126600, bc.NewHash([32]byte{0xac, 0x10, 0x41, 0x08, 0x24, 0x80, 0xe9, 0x5a, 0x9f, 0x32, 0x0a, 0x5e, 0x17, 0x7b, 0x01, 0x8d, 0x0d, 0x0d, 0x3d, 0xfc, 0xa7, 0x1d, 0x81, 0x5f, 0x13, 0xb4, 0xad, 0x0f, 0xc6, 0xde, 0x7a, 0x10})},
		{131260, bc.NewHash([32]byte{0xdf, 0x18, 0xb5, 0xb1, 0x6f, 0x5f, 0xd2, 0x77, 0x7c, 0xab, 0xb8, 0x59, 0xcb, 0x13, 0x64, 0xce, 0x06, 0x06, 0x51, 0x39, 0x89, 0x30, 0x1b, 0x69, 0xd6, 0x00, 0xec, 0xd8, 0xfa, 0xd2, 0x09, 0x93})},
		{157000, bc.NewHash([32]byte{0xb7, 0x70, 0x38, 0x4c, 0x81, 0x32, 0xaf, 0x12, 0x8d, 0xfa, 0xb4, 0xeb, 0x46, 0x4e, 0xb7, 0xeb, 0x66, 0x14, 0xd9, 0x24, 0xc2, 0xd1, 0x0c, 0x9c, 0x14, 0x20, 0xc9, 0xea, 0x0e, 0x85, 0xc8, 0xc3})},
		{180000, bc.NewHash([32]byte{0x3c, 0x2a, 0x91, 0x55, 0xf3, 0x36, 0x6a, 0x5a, 0x60, 0xcf, 0x84, 0x42, 0xec, 0x4d, 0x0c, 0x63, 0xbc, 0x34, 0xe9, 0x1d, 0x1c, 0x6b, 0xb0, 0xf0, 0x50, 0xf3, 0xfb, 0x2d, 0xf6, 0xa1, 0xd9, 0x5c})},
		{191000, bc.NewHash([32]byte{0x09, 0x4f, 0xe3, 0x23, 0x91, 0xb5, 0x11, 0x18, 0x68, 0xcc, 0x99, 0x9f, 0xeb, 0x95, 0xf9, 0xcc, 0xa5, 0x27, 0x6a, 0xf9, 0x0e, 0xda, 0x1b, 0xc6, 0x2e, 0x03, 0x29, 0xfe, 0x08, 0xdd, 0x2b, 0x01})},
		{205000, bc.NewHash([32]byte{0x6f, 0xdd, 0x87, 0x26, 0x73, 0x3f, 0x0b, 0xc7, 0x58, 0x64, 0xa4, 0xdf, 0x45, 0xe4, 0x50, 0x27, 0x68, 0x38, 0x18, 0xb9, 0xa9, 0x44, 0x56, 0x20, 0x34, 0x68, 0xd8, 0x68, 0x72, 0xdb, 0x65, 0x6f})},
		{219700, bc.NewHash([32]byte{0x98, 0x49, 0x8d, 0x4b, 0x7e, 0xe9, 0x44, 0x55, 0xc1, 0x07, 0xdd, 0x9a, 0xba, 0x6b, 0x49, 0x92, 0x61, 0x15, 0x03, 0x4f, 0x59, 0x42, 0x35, 0x74, 0xea, 0x3b, 0xdb, 0x2c, 0x53, 0x11, 0x75, 0x74})},
		{240000, bc.NewHash([32]byte{0x35, 0x16, 0x65, 0x58, 0xf4, 0xef, 0x24, 0x82, 0x43, 0xbb, 0x15, 0x79, 0xd4, 0xfe, 0x1b, 0x14, 0x9f, 0xe9, 0xf0, 0xe0, 0x48, 0x72, 0x86, 0x68, 0xa7, 0xb9, 0xda, 0x58, 0x66, 0x3b, 0x1c, 0xcb})},
		{270000, bc.NewHash([32]byte{0x9d, 0x6f, 0xcc, 0xd8, 0xb8, 0xe4, 0x8c, 0x17, 0x52, 0x9a, 0xe6, 0x1b, 0x40, 0x60, 0xe0, 0xe3, 0x6d, 0x1e, 0x89, 0xc0, 0x26, 0xdf, 0x1c, 0x28, 0x18, 0x0d, 0x29, 0x0c, 0x9b, 0x15, 0xcc, 0x97})},
		{300000, bc.NewHash([32]byte{0xa2, 0x85, 0x84, 0x6c, 0xe0, 0x3e, 0x1d, 0x68, 0x98, 0x7d, 0x93, 0x21, 0xea, 0xcc, 0x1d, 0x07, 0x88, 0xd1, 0x4c, 0x77, 0xa3, 0xd7, 0x55, 0x8a, 0x2b, 0x4a, 0xf7, 0x4d, 0x50, 0x14, 0x53, 0x5d})},
		{320000, bc.NewHash([32]byte{0xc6, 0xe7, 0x91, 0x6f, 0xcb, 0x7a, 0x42, 0x5d, 0xd6, 0x22, 0xef, 0x5d, 0x6a, 0x5c, 0xc1, 0x91, 0xa9, 0xd9, 0x06, 0x44, 0xcf, 0x36, 0x43, 0x55, 0xfe, 0x45, 0xaf, 0x24, 0x07, 0x31, 0x23, 0xc6})},
		{350000, bc.NewHash([32]byte{0xf4, 0x88, 0x6f, 0x9a, 0x17, 0x9a, 0x1c, 0x2a, 0x43, 0x9f, 0xc5, 0xae, 0x2f, 0xe4, 0xa6, 0x33, 0x71, 0xb4, 0xcd, 0x83, 0xc5, 0x23, 0xc1, 0x14, 0xb2, 0xb0, 0xa8, 0x43, 0xf2, 0xa1, 0x4b, 0x5c})},
		{380000, bc.NewHash([32]byte{0xf1, 0x00, 0x41, 0xcc, 0xea, 0xf0, 0x67, 0x98, 0x49, 0x89, 0x5f, 0xa6, 0xa0, 0x8d, 0x33, 0x04, 0x93, 0x1b, 0xf8, 0x49, 0x76, 0xe1, 0x22, 0xe3, 0xce, 0xc2, 0x6f, 0x19, 0xff, 0x4f, 0x80, 0xf1})},
		{420000, bc.NewHash([32]byte{0xdc, 0x99, 0xc7, 0x2d, 0x52, 0xa9, 0xc7, 0x2f, 0xdb, 0xcc, 0xf2, 0xb0, 0x51, 0xa2, 0x31, 0xde, 0x01, 0x06, 0x3d, 0x5d, 0x5e, 0x1a, 0x35, 0xd7, 0x1f, 0xf0, 0x5c, 0x2b, 0x59, 0x5c, 0x06, 0x32})},
		{460000, bc.NewHash([32]byte{0xf7, 0xfa, 0x12, 0xc9, 0xab, 0x6c, 0xcd, 0x32, 0x4e, 0xb0, 0x4a, 0x57, 0x6a, 0xe0, 0x9e, 0xd4, 0xa3, 0x59, 0x2d, 0x83, 0x5a, 0x27, 0xbb, 0x1b, 0xc0, 0xcd, 0x12, 0x94, 0x9e, 0xaf, 0x03, 0x7d})},
		{500000, bc.NewHash([32]byte{0x35, 0xc7, 0x54, 0x9a, 0x4b, 0x7b, 0x61, 0x2a, 0xe9, 0x4e, 0x16, 0xd6, 0x29, 0x8b, 0x97, 0x18, 0x0c, 0xa5, 0x73, 0x6a, 0xcb, 0x79, 0xfb, 0x2a, 0xc6, 0x4a, 0xf2, 0xc8, 0x2a, 0x18, 0x35, 0xe6})},
		{570000, bc.NewHash([32]byte{0x46, 0xd2, 0x28, 0xb4, 0xf3, 0xa2, 0x92, 0xde, 0x9d, 0xb3, 0x1c, 0x86, 0x4b, 0xb9, 0x94, 0x65, 0xad, 0x31, 0xa4, 0xa7, 0x28, 0x94, 0x30, 0x8c, 0xea, 0xa0, 0xaa, 0xf1, 0x9a, 0xe0, 0xbd, 0x9c})},
	},
}

// TestNetParams is the config for test-net
var TestNetParams = Params{
	Name:            "test",
	Bech32HRPSegwit: "ty",
	DefaultPort:     "46656",
	DNSSeeds:        []string{"testnetseed.eiyaro.org"},
	Checkpoints: []Checkpoint{
		{10303, bc.NewHash([32]byte{0x3e, 0x94, 0x5d, 0x35, 0x70, 0x30, 0xd4, 0x3b, 0x3d, 0xe3, 0xdd, 0x80, 0x67, 0x29, 0x9a, 0x5e, 0x09, 0xf9, 0xfb, 0x2b, 0xad, 0x5f, 0x92, 0xc8, 0x69, 0xd1, 0x42, 0x39, 0x74, 0x9a, 0xd1, 0x1c})},
		{40000, bc.NewHash([32]byte{0x6b, 0x13, 0x9a, 0x5b, 0x76, 0x77, 0x9b, 0xd4, 0x1c, 0xec, 0x53, 0x68, 0x44, 0xbf, 0xf4, 0x48, 0x94, 0x3d, 0x16, 0xe3, 0x9b, 0x2e, 0xe8, 0xa1, 0x0f, 0xa0, 0xbc, 0x7d, 0x2b, 0x17, 0x55, 0xfc})},
		{78000, bc.NewHash([32]byte{0xa9, 0x03, 0xc0, 0x0c, 0x62, 0x1a, 0x3d, 0x00, 0x7f, 0xd8, 0x5d, 0x51, 0xba, 0x43, 0xe4, 0xd0, 0xe3, 0xc5, 0xd4, 0x8f, 0x30, 0xb5, 0x5f, 0xa5, 0x77, 0x62, 0xd8, 0x8b, 0x11, 0x81, 0x5f, 0xb4})},
		{82000, bc.NewHash([32]byte{0x56, 0xb1, 0xba, 0x23, 0x69, 0x5c, 0x8f, 0x51, 0x4e, 0x23, 0xc0, 0xae, 0xaa, 0x25, 0x08, 0xc5, 0x85, 0xf3, 0x7c, 0xd1, 0xc6, 0x15, 0xa2, 0x51, 0xda, 0x79, 0x4f, 0x08, 0x13, 0x66, 0xc9, 0x85})},
		{83200, bc.NewHash([32]byte{0xb4, 0x6f, 0xc5, 0xcf, 0xa3, 0x3d, 0xe1, 0x11, 0x71, 0x68, 0x40, 0x68, 0x0c, 0xe7, 0x4c, 0xaf, 0x5a, 0x11, 0xfe, 0x82, 0xbc, 0x36, 0x88, 0x0f, 0xbd, 0x04, 0xf0, 0xc4, 0x86, 0xd4, 0xd6, 0xd5})},
		{93000, bc.NewHash([32]byte{0x6f, 0x4f, 0x37, 0x5f, 0xe9, 0xfb, 0xdf, 0x66, 0x60, 0x0e, 0xf0, 0x39, 0xb7, 0x18, 0x26, 0x75, 0xa0, 0x9a, 0xa5, 0x9b, 0x83, 0xc9, 0x9a, 0x25, 0x45, 0xb8, 0x7d, 0xd4, 0x99, 0x24, 0xa2, 0x8a})},
		{113300, bc.NewHash([32]byte{0x7a, 0x69, 0x75, 0xa5, 0xf6, 0xb6, 0x94, 0xf3, 0x94, 0xa2, 0x63, 0x91, 0x28, 0xb6, 0xab, 0x7e, 0xf9, 0x71, 0x27, 0x5a, 0xe2, 0x59, 0xd3, 0xff, 0x70, 0x6e, 0xcb, 0xd8, 0xd8, 0x30, 0x9c, 0xc4})},
		{235157, bc.NewHash([32]byte{0xfa, 0x76, 0x36, 0x3e, 0x9e, 0x58, 0xea, 0xe4, 0x7d, 0x26, 0x70, 0x7e, 0xf3, 0x8b, 0xfd, 0xad, 0x1a, 0x99, 0xf7, 0x4c, 0xac, 0xc6, 0x80, 0x99, 0x58, 0x10, 0x13, 0x66, 0x4b, 0x8c, 0x39, 0x4f})},
		{252383, bc.NewHash([32]byte{0xa1, 0xaa, 0xe6, 0xd9, 0x42, 0x94, 0x99, 0x7b, 0x9b, 0x71, 0x5b, 0xf5, 0x23, 0x9a, 0xee, 0x92, 0x27, 0x84, 0x4c, 0x32, 0x47, 0xf2, 0xf2, 0xd9, 0xe3, 0xd7, 0x6a, 0x2b, 0xbe, 0xc3, 0x1f, 0x50})},
		{320000, bc.NewHash([32]byte{0xb5, 0x9f, 0xeb, 0x44, 0x4a, 0xfd, 0x67, 0x2a, 0x1c, 0x99, 0x8e, 0xc8, 0x48, 0xac, 0xeb, 0xfe, 0x80, 0xf6, 0x58, 0x7a, 0xf1, 0x7d, 0xea, 0xd6, 0xc7, 0xe5, 0x93, 0x3f, 0xae, 0x65, 0x88, 0xe6})},
		{360000, bc.NewHash([32]byte{0xad, 0x30, 0x55, 0x89, 0xf6, 0xd3, 0x7c, 0x13, 0x81, 0x8d, 0xeb, 0x09, 0xce, 0xea, 0x08, 0x64, 0xc4, 0x39, 0xd5, 0x80, 0x61, 0x88, 0xe2, 0xa5, 0xc8, 0x8f, 0xd7, 0x64, 0x4b, 0x0b, 0x5a, 0x6b})},
		{400000, bc.NewHash([32]byte{0xc2, 0xa5, 0xa3, 0x64, 0x51, 0x01, 0xe3, 0x1e, 0x57, 0xfb, 0x48, 0x15, 0x41, 0xcf, 0xaa, 0x67, 0x66, 0xd1, 0xb3, 0x2d, 0x91, 0xd2, 0x3d, 0xef, 0x53, 0xcd, 0x59, 0x97, 0x55, 0x04, 0xab, 0xc8})},
		{440000, bc.NewHash([32]byte{0xca, 0x56, 0xff, 0xf7, 0x94, 0x3c, 0x1b, 0xc7, 0x00, 0x52, 0x50, 0xa3, 0x3f, 0x75, 0xa3, 0x4b, 0xeb, 0x06, 0xc8, 0xba, 0xe3, 0x45, 0x9f, 0x93, 0x2a, 0x8c, 0x6c, 0xa3, 0x85, 0x33, 0x2b, 0x82})},
	},
}

// SoloNetParams is the config for test-net
var SoloNetParams = Params{
	Name:            "solo",
	Bech32HRPSegwit: "sy",
	Checkpoints:     []Checkpoint{},
}
