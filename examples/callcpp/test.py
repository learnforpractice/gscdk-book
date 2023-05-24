import os
import sys
import json
import struct
import pytest

test_dir = os.path.dirname(__file__)
sys.path.append(os.path.join(test_dir, '..'))

from ipyeos import log
from ipyeos import chaintester
chaintester.chain_config['contracts_console'] = True

logger = log.get_logger(__name__)

def update_auth(chain, account):
    a = {
        "account": account,
        "permission": "active",
        "parent": "owner",
        "auth": {
            "threshold": 1,
            "keys": [
                {
                    "key": 'EOS6AjF6hvF7GSuSd4sCgfPKq5uWaXvGM2aQtEUCwmEHygQaqxBSV',
                    "weight": 1
                }
            ],
            "accounts": [{"permission":{"actor":account,"permission": 'eosio.code'}, "weight":1}],
            "waits": []
        }
    }
    chain.push_action('eosio', 'updateauth', a, {account:'active'})

def init_tester():
    chain = chaintester.ChainTester()
    update_auth(chain, 'hello')
    return chain

def chain_test(fn):
    def call():
        chain = init_tester()
        ret = fn(chain)
        chain.free()
        return ret
    return call

class NewChainTester():
    def __init__(self):
        self.tester = None

    def __enter__(self):
        self.tester = init_tester()
        return self.tester

    def __exit__(self, type, value, traceback):
        self.tester.free()

test_dir = os.path.dirname(__file__)
def deploy_contract(tester, package_name):
    with open(f'{test_dir}/{package_name}.wasm', 'rb') as f:
        code = f.read()
    with open(f'{test_dir}/{package_name}.abi', 'rb') as f:
        abi = f.read()
    tester.deploy_contract('hello', code, abi)

@chain_test
def test_say_hello(tester):
    deploy_contract(tester, 'callcpp')
    r = tester.push_action('hello', 'sayhello', b'', {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
