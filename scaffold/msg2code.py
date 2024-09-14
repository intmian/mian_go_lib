import os
import re
import sys

# 定义正则表达式匹配目标代码
cmd_pattern = r'const\s+Cmd(\w+)\s+.*Cmd\s+=\s+"(\w+)"'
req_pattern = r'type\s+(\w+)Req\s+struct\s+{'
ret_pattern = r'type\s+(\w+)Ret\s+struct\s+{'

content = ""
addr = ""
# 读取文件内容,如果参数含有文件名，打开文件
if len(sys.argv) > 1 and os.path.exists(sys.argv[1]):
    addr = sys.argv[1]
else:
    addr = "msg.go"
with open(sys.argv[1], "r", encoding="utf-8") as file:
    content = file.read()

# 匹配CmdXxx, XxxReq, XxxRet
cmd_matches = re.findall(cmd_pattern, content)
req_matches = re.findall(req_pattern, content)
ret_matches = re.findall(ret_pattern, content)
print(cmd_matches, req_matches, ret_matches)

# 确保XxxReq和XxxRet匹配一致
if len(req_matches) != len(ret_matches) or len(cmd_matches) != len(req_matches):
    print("匹配失败")
    exit()

# 输出生成的代码
print('switch msg.Cmd(){')
for i in range(len(cmd_matches)):
    cmd_name, cmd_value = cmd_matches[i]
    req_name = req_matches[i]
    ret_name = ret_matches[i]
    
    # case Xxx 部分
    print(f'case Cmd{cmd_name}:')
    print(f'    return backendshare.HandleRpcTool("{cmd_value}", msg, valid, s.On{cmd_name})')
print('}')
    
# 输出生成的代码
for i in range(len(cmd_matches)):
    cmd_name, cmd_value = cmd_matches[i]
    req_name = req_matches[i]
    ret_name = ret_matches[i]
    
    # OnXxx 函数部分
    print(f'func (s *Service) On{cmd_name}(valid backendshare.Valid, req {req_name}Req) (ret {ret_name}Ret, err error) {{')
    print(f'    // TODO')
    print(f'}}\n')
    
# 输出生成的代码
for i in range(len(cmd_matches)):
    cmd_name, cmd_value = cmd_matches[i]
    req_name = req_matches[i]
    ret_name = ret_matches[i]
    
    # sendXxx 函数部分
    print(f'export function send{cmd_name}(req, callback) {{')
    print(f'    UniPost(config.api_base_url + \'/{cmd_value}\', req).then(callback)')
    print(f'}}\n')