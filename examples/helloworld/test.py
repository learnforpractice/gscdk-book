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

def init_tester():
    chain = chaintester.ChainTester()
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
def test_helloworld(tester):
    deploy_contract(tester, 'helloworld')

    r = tester.push_action('hello', 'sayhello', {}, {'hello': 'active'})
    logger.info('++++++elapsed: %s', r['elapsed'])
    tester.produce_block()
