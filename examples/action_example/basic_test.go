package main

import (
	"context"
	"os"
	"testing"

	"github.com/uuosio/chaintester"
)

var ctx = context.Background()

func initTest() *chaintester.ChainTester {
	tester := chaintester.NewChainTester()
	testCoverage := os.Getenv("TEST_COVERAGE")
	if testCoverage == "TRUE" || testCoverage == "true" || testCoverage == "1" {
		tester.SetNativeApply("hello", ContractApply)
	}
	return tester
}

func TestHello(t *testing.T) {
	permissions := `
	{
		"hello": "active"
	}
	`

	tester := initTest()
	defer tester.FreeChain()

	updateAuthArgs := `{
		"account": "hello",
		"permission": "active",
		"parent": "owner",
		"auth": {
			"threshold": 1,
			"keys": [
				{
					"key": "EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV",
					"weight": 1
				}
			],
			"accounts": [{"permission":{"actor": "hello", "permission": "eosio.code"}, "weight":1}],
			"waits": []
		}
	}`
	tester.PushAction("eosio", "updateauth", updateAuthArgs, permissions)

	err := tester.DeployContract("hello", "action_example.wasm", "action_example.abi")
	if err != nil {
		panic(err)
	}
	tester.ProduceBlock()

	oldBalance := tester.GetBalance("hello")
	t.Logf("++++++++old Balance: %v", oldBalance)

	_, err = tester.PushAction("hello", "testaction1", "", permissions)
	if err != nil {
		panic(err)
	}

	newBalance := tester.GetBalance("hello")
	t.Logf("++++++++new balance: %v", newBalance)
	if oldBalance-newBalance != 10000 {
		panic("invalid balance")
	}

	oldBalance = tester.GetBalance("hello")
	t.Logf("++++++++old Balance: %v", oldBalance)

	_, err = tester.PushAction("hello", "testaction2", "", permissions)
	if err != nil {
		panic(err)
	}

	newBalance = tester.GetBalance("hello")
	t.Logf("++++++++new balance: %v", newBalance)
	if oldBalance-newBalance != 10000 {
		panic("invalid balance")
	}
}
