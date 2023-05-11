// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package OracleV2

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// OracleV2MetaData contains all meta data concerning the OracleV2 contract.
var OracleV2MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"timestamp\",\"type\":\"uint128\"}],\"name\":\"OracleUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newUpdater\",\"type\":\"address\"}],\"name\":\"UpdaterAddressChange\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"getValue\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"},{\"internalType\":\"uint128\",\"name\":\"value\",\"type\":\"uint128\"},{\"internalType\":\"uint128\",\"name\":\"timestamp\",\"type\":\"uint128\"}],\"name\":\"setValue\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOracleUpdaterAddress\",\"type\":\"address\"}],\"name\":\"updateOracleUpdaterAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"values\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// OracleV2ABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleV2MetaData.ABI instead.
var OracleV2ABI = OracleV2MetaData.ABI

// OracleV2 is an auto generated Go binding around an Ethereum contract.
type OracleV2 struct {
	OracleV2Caller     // Read-only binding to the contract
	OracleV2Transactor // Write-only binding to the contract
	OracleV2Filterer   // Log filterer for contract events
}

// OracleV2Caller is an auto generated read-only Go binding around an Ethereum contract.
type OracleV2Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleV2Transactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleV2Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleV2Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleV2Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleV2Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleV2Session struct {
	Contract     *OracleV2         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleV2CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleV2CallerSession struct {
	Contract *OracleV2Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// OracleV2TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleV2TransactorSession struct {
	Contract     *OracleV2Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// OracleV2Raw is an auto generated low-level Go binding around an Ethereum contract.
type OracleV2Raw struct {
	Contract *OracleV2 // Generic contract binding to access the raw methods on
}

// OracleV2CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleV2CallerRaw struct {
	Contract *OracleV2Caller // Generic read-only contract binding to access the raw methods on
}

// OracleV2TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleV2TransactorRaw struct {
	Contract *OracleV2Transactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleV2 creates a new instance of OracleV2, bound to a specific deployed contract.
func NewOracleV2(address common.Address, backend bind.ContractBackend) (*OracleV2, error) {
	contract, err := bindOracleV2(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleV2{OracleV2Caller: OracleV2Caller{contract: contract}, OracleV2Transactor: OracleV2Transactor{contract: contract}, OracleV2Filterer: OracleV2Filterer{contract: contract}}, nil
}

// NewOracleV2Caller creates a new read-only instance of OracleV2, bound to a specific deployed contract.
func NewOracleV2Caller(address common.Address, caller bind.ContractCaller) (*OracleV2Caller, error) {
	contract, err := bindOracleV2(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleV2Caller{contract: contract}, nil
}

// NewOracleV2Transactor creates a new write-only instance of OracleV2, bound to a specific deployed contract.
func NewOracleV2Transactor(address common.Address, transactor bind.ContractTransactor) (*OracleV2Transactor, error) {
	contract, err := bindOracleV2(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleV2Transactor{contract: contract}, nil
}

// NewOracleV2Filterer creates a new log filterer instance of OracleV2, bound to a specific deployed contract.
func NewOracleV2Filterer(address common.Address, filterer bind.ContractFilterer) (*OracleV2Filterer, error) {
	contract, err := bindOracleV2(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleV2Filterer{contract: contract}, nil
}

// bindOracleV2 binds a generic wrapper to an already deployed contract.
func bindOracleV2(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleV2MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleV2 *OracleV2Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleV2.Contract.OracleV2Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleV2 *OracleV2Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleV2.Contract.OracleV2Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleV2 *OracleV2Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleV2.Contract.OracleV2Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleV2 *OracleV2CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _OracleV2.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleV2 *OracleV2TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleV2.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleV2 *OracleV2TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleV2.Contract.contract.Transact(opts, method, params...)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_OracleV2 *OracleV2Caller) GetValue(opts *bind.CallOpts, key string) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _OracleV2.contract.Call(opts, &out, "getValue", key)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_OracleV2 *OracleV2Session) GetValue(key string) (*big.Int, *big.Int, error) {
	return _OracleV2.Contract.GetValue(&_OracleV2.CallOpts, key)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_OracleV2 *OracleV2CallerSession) GetValue(key string) (*big.Int, *big.Int, error) {
	return _OracleV2.Contract.GetValue(&_OracleV2.CallOpts, key)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_OracleV2 *OracleV2Caller) Values(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _OracleV2.contract.Call(opts, &out, "values", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_OracleV2 *OracleV2Session) Values(arg0 string) (*big.Int, error) {
	return _OracleV2.Contract.Values(&_OracleV2.CallOpts, arg0)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_OracleV2 *OracleV2CallerSession) Values(arg0 string) (*big.Int, error) {
	return _OracleV2.Contract.Values(&_OracleV2.CallOpts, arg0)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_OracleV2 *OracleV2Transactor) SetValue(opts *bind.TransactOpts, key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _OracleV2.contract.Transact(opts, "setValue", key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_OracleV2 *OracleV2Session) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _OracleV2.Contract.SetValue(&_OracleV2.TransactOpts, key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_OracleV2 *OracleV2TransactorSession) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _OracleV2.Contract.SetValue(&_OracleV2.TransactOpts, key, value, timestamp)
}

// UpdateOracleUpdaterAddress is a paid mutator transaction binding the contract method 0x6aa45efc.
//
// Solidity: function updateOracleUpdaterAddress(address newOracleUpdaterAddress) returns()
func (_OracleV2 *OracleV2Transactor) UpdateOracleUpdaterAddress(opts *bind.TransactOpts, newOracleUpdaterAddress common.Address) (*types.Transaction, error) {
	return _OracleV2.contract.Transact(opts, "updateOracleUpdaterAddress", newOracleUpdaterAddress)
}

// UpdateOracleUpdaterAddress is a paid mutator transaction binding the contract method 0x6aa45efc.
//
// Solidity: function updateOracleUpdaterAddress(address newOracleUpdaterAddress) returns()
func (_OracleV2 *OracleV2Session) UpdateOracleUpdaterAddress(newOracleUpdaterAddress common.Address) (*types.Transaction, error) {
	return _OracleV2.Contract.UpdateOracleUpdaterAddress(&_OracleV2.TransactOpts, newOracleUpdaterAddress)
}

// UpdateOracleUpdaterAddress is a paid mutator transaction binding the contract method 0x6aa45efc.
//
// Solidity: function updateOracleUpdaterAddress(address newOracleUpdaterAddress) returns()
func (_OracleV2 *OracleV2TransactorSession) UpdateOracleUpdaterAddress(newOracleUpdaterAddress common.Address) (*types.Transaction, error) {
	return _OracleV2.Contract.UpdateOracleUpdaterAddress(&_OracleV2.TransactOpts, newOracleUpdaterAddress)
}

// OracleV2OracleUpdateIterator is returned from FilterOracleUpdate and is used to iterate over the raw logs and unpacked data for OracleUpdate events raised by the OracleV2 contract.
type OracleV2OracleUpdateIterator struct {
	Event *OracleV2OracleUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleV2OracleUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleV2OracleUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleV2OracleUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleV2OracleUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleV2OracleUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleV2OracleUpdate represents a OracleUpdate event raised by the OracleV2 contract.
type OracleV2OracleUpdate struct {
	Key       string
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdate is a free log retrieval operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_OracleV2 *OracleV2Filterer) FilterOracleUpdate(opts *bind.FilterOpts) (*OracleV2OracleUpdateIterator, error) {

	logs, sub, err := _OracleV2.contract.FilterLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return &OracleV2OracleUpdateIterator{contract: _OracleV2.contract, event: "OracleUpdate", logs: logs, sub: sub}, nil
}

// WatchOracleUpdate is a free log subscription operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_OracleV2 *OracleV2Filterer) WatchOracleUpdate(opts *bind.WatchOpts, sink chan<- *OracleV2OracleUpdate) (event.Subscription, error) {

	logs, sub, err := _OracleV2.contract.WatchLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleV2OracleUpdate)
				if err := _OracleV2.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOracleUpdate is a log parse operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_OracleV2 *OracleV2Filterer) ParseOracleUpdate(log types.Log) (*OracleV2OracleUpdate, error) {
	event := new(OracleV2OracleUpdate)
	if err := _OracleV2.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleV2UpdaterAddressChangeIterator is returned from FilterUpdaterAddressChange and is used to iterate over the raw logs and unpacked data for UpdaterAddressChange events raised by the OracleV2 contract.
type OracleV2UpdaterAddressChangeIterator struct {
	Event *OracleV2UpdaterAddressChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleV2UpdaterAddressChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleV2UpdaterAddressChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleV2UpdaterAddressChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleV2UpdaterAddressChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleV2UpdaterAddressChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleV2UpdaterAddressChange represents a UpdaterAddressChange event raised by the OracleV2 contract.
type OracleV2UpdaterAddressChange struct {
	NewUpdater common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdaterAddressChange is a free log retrieval operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_OracleV2 *OracleV2Filterer) FilterUpdaterAddressChange(opts *bind.FilterOpts) (*OracleV2UpdaterAddressChangeIterator, error) {

	logs, sub, err := _OracleV2.contract.FilterLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return &OracleV2UpdaterAddressChangeIterator{contract: _OracleV2.contract, event: "UpdaterAddressChange", logs: logs, sub: sub}, nil
}

// WatchUpdaterAddressChange is a free log subscription operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_OracleV2 *OracleV2Filterer) WatchUpdaterAddressChange(opts *bind.WatchOpts, sink chan<- *OracleV2UpdaterAddressChange) (event.Subscription, error) {

	logs, sub, err := _OracleV2.contract.WatchLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleV2UpdaterAddressChange)
				if err := _OracleV2.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdaterAddressChange is a log parse operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_OracleV2 *OracleV2Filterer) ParseUpdaterAddressChange(log types.Log) (*OracleV2UpdaterAddressChange, error) {
	event := new(OracleV2UpdaterAddressChange)
	if err := _OracleV2.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
