# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

import signal
import time
from functools import total_ordering
from threading import Lock
from threading import Thread

import requests
from constants import INDEXING_DISABLED_ERROR

from github import get_chain_config
from helpers import get_chains_info
from redis_functions import redisSetEndpoint
from redis_functions import setPrimaryEndpoint
from redis_functions import setSecondaryEndpoint
from redis_functions import setTertiaryEndpoint
from redis_functions import flushChains


@total_ordering
class Endpoint:

    def __init__(self, url: str, height: int, latency: float):
        self.url = url
        self.height = height
        self.latency = latency

    def __eq__(self, other):
        return self.url == other.url

    def __lt__(self, other):
        if self.height == -1:
            return True

        if self.height != other.height:
            return self.height < other.height

        if self.latency == -1:
            return True

        return self.latency > other.latency

    def __repr__(self):
        return f'Endpoint: {self.url} - Height: {self.height} - Latency: {self.latency}'


class EndpointSafeList:

    def __init__(self):
        self.lock = Lock()
        self.elements: list[Endpoint] = []

    def reset(self):
        with self.lock:
            self.elements = []

    def add_element(self, url: str, height: int, latency: float):
        with self.lock:
            self.elements.append(Endpoint(url, height, latency))


def ping_jrpc(endpoint: str, jrpc_list: EndpointSafeList) -> None:
    url = f'{endpoint}/status'
    transaction_url = f'{endpoint}/tx?hash=0x0000000000000000000000000000000000000000000000000000000000000000'
                                           
    try:
        transaction_res = requests.get(transaction_url, timeout=5)
        transaction_parsed = transaction_res.json()
        indexing_disabled = transaction_parsed.get('error') and INDEXING_DISABLED_ERROR in transaction_parsed.get('error', {}).get('data')
        if not indexing_disabled:
            res = requests.get(url, timeout=1)
            parsed = res.json()
            height = int(parsed['result']['sync_info']['latest_block_height'])
            latency = res.elapsed.total_seconds()
            # We need tx_index to check the transaction status
            tx_index = False
            try:
                if parsed['result']['node_info']['other']['tx_index'] == 'on':
                    tx_index = True
            except Exception:
                pass  
            if tx_index is True:
                jrpc_list.add_element(endpoint, height, latency)

    except Exception:
        jrpc_list.add_element(endpoint, -1, -1)


def ping_web3(endpoint: str, web3_list: EndpointSafeList) -> None:
    url = f'{endpoint}/status'
    try:
        res = requests.post(
            url, json='{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}', timeout=1)
        parsed = res.json()
        height = int(parsed['result'], base=16)
        latency = res.elapsed.total_seconds()
        web3_list.add_element(endpoint, height, latency)
    except Exception:
        web3_list.add_element(endpoint, -1, -1)


def ping_rest(endpoint: str, rest_list: EndpointSafeList, tendermint_exposed=True) -> None:
    if tendermint_exposed:
        url = f'{endpoint}/cosmos/base/tendermint/v1beta1/blocks/latest'
        try:
            res = requests.get(url, timeout=1)
            parsed = res.json()
            height = parsed.get('block', {}).get('header', {}).get('height', None)
            if  height is not None:
                height = int(parsed['block']['header']['height'])
                latency = res.elapsed.total_seconds()
                rest_list.add_element(endpoint, height, latency)
            else:
                rest_list.add_element(endpoint, -1, -1)
        except Exception:
            rest_list.add_element(endpoint, -1, -1)
    else:
        # If the chain rest api doesn't have the /base/tendermint endpoints exposed
        # We are going to just compare using the latency
        try:
            res = requests.get(
                f'{endpoint}/cosmos/auth/v1beta1/params', timeout=1)
            rest_list.add_element(endpoint, 0, res.elapsed.total_seconds())
        except Exception:
            rest_list.add_element(endpoint, -1, -1)


def process_chain(chain: str, chain_info):
    tendermint_exposed = chain != 'GRAVITYBRIDGE'
    rest_threads = []
    rest_list = EndpointSafeList()
    for rest in chain_info[chain]['rest']:
        t = Thread(target=ping_rest, args=(
            rest,
            rest_list,
            tendermint_exposed,
        ))
        t.start()
        rest_threads.append(t)

    jrpc_threads = []
    jrpc_list = EndpointSafeList()
    for jrpc in chain_info[chain]['jrpc']:
        t = Thread(target=ping_jrpc, args=(
            jrpc,
            jrpc_list,
        ))
        t.start()
        jrpc_threads.append(t)

    web3_list = EndpointSafeList()
    web3_threads = []
    if 'web3' in chain_info[chain]:
        for web3 in chain_info[chain]['web3']:
            t = Thread(target=ping_web3, args=(
                web3,
                web3_list,
            ))
            t.start()
            web3_threads.append(t)

    for t in rest_threads:
        t.join()

    for t in jrpc_threads:
        t.join()

    for t in web3_threads:
        t.join()

    if len(rest_list.elements) > 2:
        rest_list.elements.sort(reverse=True)
        setPrimaryEndpoint(chain, 'rest', rest_list.elements[0].url)
        setSecondaryEndpoint(chain, 'rest', rest_list.elements[1].url)
        setTertiaryEndpoint(chain, 'rest', rest_list.elements[2].url)

    if len(jrpc_list.elements) > 2:
        jrpc_list.elements.sort(reverse=True)
        setPrimaryEndpoint(chain, 'jrpc', jrpc_list.elements[0].url)
        setSecondaryEndpoint(chain, 'jrpc', jrpc_list.elements[1].url)
        setTertiaryEndpoint(chain, 'jrpc', jrpc_list.elements[2].url)

    if len(web3_list.elements) > 2:
        web3_list.elements.sort(reverse=True)
        redisSetEndpoint(chain, 'web3', 0, 'https://evmos-evm.publicnode.com')
        setPrimaryEndpoint(chain, 'web3', web3_list.elements[0].url)
        setSecondaryEndpoint(chain, 'web3', web3_list.elements[1].url)
        setTertiaryEndpoint(chain, 'web3', web3_list.elements[2].url)


running = True


def main():
    global running
    threads = []
    attempt = 0
    while running:
        try:
            print('Getting chain config...')
            chain_data = get_chain_config()
            chain_info = get_chains_info(chain_data)
            start_time = time.time()
            for chain in chain_info:
                t = Thread(target=process_chain, args=(chain, chain_info, ))
                t.start()
                threads.append(t)

            for t in threads:
                t.join()

            print(f'Time used: {time.time() - start_time}')
            threads = []
            attempt = 0
            time.sleep(5)
            
        except Exception as e:
            print('Failed to get chain config, flushing redis and trying again')
            print(e)
            attempt += 1
            flushChains()
            if attempt > 5:
                time.sleep(5)

               



def signal_handler(sig, frame):
    global running
    _ = sig
    _ = frame
    print('You pressed Ctrl+C!')
    running = False


signal.signal(signal.SIGINT, signal_handler)

if __name__ == '__main__':
    main()
