"""
spider for ** category data
"""

import io
import time
import json
import requests

url = "https://upload.taobao.com/auction/json/reload_cats.htm"

querystring = {"customId": "", "fenxiaoProduct": ""}

headers = {}


def parseSub(sid):
    print("parsing sid: ", sid)
    # request
    payload = {"sid": sid, "path": "next"}
    response = requests.request(
        "POST", url, data=payload, headers=headers, params=querystring)
    data = response.json()
    # print(data, "\n\n")

    if not data:
        return

    # parse
    for node in data[0]["data"]:
        if not node["data"]:
            continue

        for cat in node["data"]:
            cat["pid"] = sid
            # categories.append(node)
            f.write(unicode(json.dumps(cat), encoding='utf8'))
            f.write(unicode("\n", encoding='utf8'))
            if cat["leaf"] == 0:
                time.sleep(0.05)
                parseSub(cat["sid"])


if __name__ == "__main__":
    with io.open("../data/categories.json", "w", encoding='utf8') as f:
        parseSub(0)
