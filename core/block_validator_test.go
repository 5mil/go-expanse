// Copyright 2015 The go-expanse Authors
// This file is part of the go-expanse library.
//
// The go-expanse library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-expanse library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-expanse library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/expanse-project/go-expanse/common"
	"github.com/expanse-project/go-expanse/core/state"
	"github.com/expanse-project/go-expanse/core/types"
	"github.com/expanse-project/go-expanse/core/vm"
	"github.com/expanse-project/go-expanse/ethdb"
	"github.com/expanse-project/go-expanse/event"
	"github.com/expanse-project/go-expanse/pow/ezp"
)

func proc() (Validator, *BlockChain) {
	db, _ := ethdb.NewMemDatabase()
	var mux event.TypeMux

	WriteTestNetGenesisBlock(db, 0)
	blockchain, err := NewBlockChain(db, thePow(), &mux)
	if err != nil {
		fmt.Println(err)
	}
	return blockchain.validator, blockchain
}

func TestNumber(t *testing.T) {
	pow := ezp.New()
	_, chain := proc()

	statedb, _ := state.New(chain.Genesis().Root(), chain.chainDb)
	header := makeHeader(chain.Genesis(), statedb)
	header.Number = big.NewInt(3)
	err := ValidateHeader(pow, header, chain.Genesis().Header(), false, false)
	if err != BlockNumberErr {
		t.Errorf("expected block number error, got %q", err)
	}

	header = makeHeader(chain.Genesis(), statedb)

	err = ValidateHeader(pow, header, chain.Genesis().Header(), false, false)
	if err == BlockNumberErr {
		t.Errorf("didn't expect block number error")
	}
}

func TestPutReceipt(t *testing.T) {
	db, _ := ethdb.NewMemDatabase()

	var addr common.Address
	addr[0] = 1
	var hash common.Hash
	hash[0] = 2

	receipt := new(types.Receipt)
	receipt.Logs = vm.Logs{&vm.Log{
		Address:     addr,
		Topics:      []common.Hash{hash},
		Data:        []byte("hi"),
		BlockNumber: 42,
		TxHash:      hash,
		TxIndex:     0,
		BlockHash:   hash,
		Index:       0,
	}}

	WriteReceipts(db, types.Receipts{receipt})
	receipt = GetReceipt(db, common.Hash{})
	if receipt == nil {
		t.Error("expected to get 1 receipt, got none.")
	}
}
