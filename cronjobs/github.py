# Copyright Tharsis Labs Ltd.(Evmos)
# SPDX-License-Identifier:ENCL-1.0(https://github.com/evmos/backend/blob/main/LICENSE)

import requests
import base64
import json
import os
from redis_functions import getChains, getTokens, setChains, setTokens


ENVIRONMENT = os.getenv("ENVIRONMENT")

CHAINS_DIRECTORY_URL = "https://api.github.com/repos/evmos/chain-token-registry/git/trees/production?recursive=1" if ENVIRONMENT == "production" else "https://api.github.com/repos/evmos/chain-token-registry/git/trees/main?recursive=1"

PAT = os.getenv("GITHUB_KEY")
if not PAT:
    raise Exception("GITHUB_KEY environment variable must be defined")


def get_git_directory_content(path: str):
    content = []
    headers = {"authorization": "token {0}".format(PAT)}
    try:
        res = requests.get(CHAINS_DIRECTORY_URL, headers=headers, timeout=5)
        parsed = res.json()
        for entry in parsed.get('tree', []):
            if (path in entry.get('path') and 'tree' not in entry.get('type')):
                res = requests.get(entry.get('url'), headers=headers, timeout=5)
                token_detail = res.json()
                file_content = token_detail.get('content')
                if file_content is not None:
                    decoded_token_detail = base64.b64decode(file_content)
                    json_token_detail = json.loads(decoded_token_detail)
                    content.append(json_token_detail)
    except Exception as e:
        print(e)
    return content


def get_chain_config():
    data = None
    data = getChains()
    if data is None or len(data) <= 0:
        data = get_git_directory_content('chainConfig')
        setChains(data)
    return data


def get_tokens():
    data = None
    data = getTokens()
    if data is None:
        data = get_git_directory_content('tokens')
        setTokens(data)
    return data
