package app

import (
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
)

// This file config instance cost and complile cost
const (
	// DefaultNibiruInstanceCost is initially set the same as in wasmd
	DefaultJunoInstanceCost uint64 = 60_000
	// DefaultNibiruCompileCost set to a large number for testing
	DefaultJunoCompileCost uint64 = 100
)

// NibiruGasRegisterConfig is defaults plus a custom compile amount
func NibiruGasRegisterConfig() wasmkeeper.WasmGasRegisterConfig {
	gasConfig := wasmkeeper.DefaultGasRegisterConfig()
	gasConfig.InstanceCost = DefaultJunoInstanceCost
	gasConfig.CompileCost = DefaultJunoCompileCost

	return gasConfig
}

func NewNibiruWasmGasRegister() wasmkeeper.WasmGasRegister {
	return wasmkeeper.NewWasmGasRegister(NibiruGasRegisterConfig())
}
