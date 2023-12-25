# 测试用例子

import sys
import json
# 读取参数
addr = sys.argv[1]
with open(addr, 'r') as f:
    data = json.load(f)

for row in data["Rows"]:
    data = row["Data"]
    if row != None:
        print("ok")
        exit(0)
