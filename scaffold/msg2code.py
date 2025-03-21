import os
import re
import sys

msg_pattern = r'const\s+Cmd(\w+)\s+.*Cmd\s+=\s+"(\w+)"\n*type\s+(\w+)Req\s+struct\s+\{([\s\S]*?)\}\n*type\s+(\w+)Ret\s+struct\s+\{([\s\S]*?)\}'
struct_value_pattern = r'(\w+)\s*([\[\]\w.]+)\s*(\/\/.*)?\n'

class Msg:
    def __init__(self,svrName,cmd,reqStructStr,retStructStr) -> None:
        # 将cmd，分别以首字母大写和小写的形式保存
        self.firstUpperCmd = cmd
        self.firstLowerCmd = cmd[0].lower() + cmd[1:]
        # 从structStr中解析出字段名和字段类型
        self.reqStruct = self.parseStruct(reqStructStr)
        self.retStruct = self.parseStruct(retStructStr)
        
        self.svrName = svrName
        
    def parseStruct(self,structStr):
        matches = re.findall(struct_value_pattern,structStr)
        struct = []
        for match in matches:
            if len(match) >= 2:
                name = match[0]
                typ = match[1]
                isArr = False
                # 如果是[]开头，转换为[]结尾
                if typ.startswith('[]'):
                    typ = typ[2:]
                    isArr = True
                # 替换
                numberType = ['int32','int64','float32','float64','uint32','uint64','int','uint']
                if typ in numberType:
                    typ = 'number'
                elif typ == 'bool':
                    typ = 'boolean'
                if len (match) == 3:
                    typ += match[2]
                if isArr:
                    typ = typ + '[] | null'
                struct.append((name,typ))
        return struct
    
    def makeCase(self):
        case = f'case Cmd{self.firstUpperCmd}:\n'
        case += f'    return backendshare.HandleRpcTool("{self.firstLowerCmd}", msg, valid, s.On{self.firstUpperCmd})\n'
        return case
    
    def makeOnFunction(self):
        onFunction = f'func (s *Service) On{self.firstUpperCmd}(valid backendshare.Valid, req {self.firstUpperCmd}Req) (ret {self.firstUpperCmd}Ret, err error) {{\n'
        onFunction += f'    // TODO\n'
        onFunction += f'}}\n'
        onFunction += '\n'
        return onFunction
    
    def makeTypeScriptInterface(self):
        interface = ''
        if len(self.reqStruct) == 0:
            interface += f'export type {self.firstUpperCmd}Req = object\n'
        else:
            interface += f'export interface {self.firstUpperCmd}Req' + ' {\n'
            for field in self.reqStruct:
                interface += f'    {field[0]}: {field[1]}\n'
            interface += '}\n'
            interface += '\n'
        
        if len(self.retStruct) == 0:
            interface += f'export type {self.firstUpperCmd}Ret = object\n'
        else:
            interface += f'export interface {self.firstUpperCmd}Ret' + ' {\n'
            for field in self.retStruct:
                interface += f'    {field[0]}: {field[1]}\n'
            interface += '}\n'
        interface += '\n'
        interface += '\n'
        return interface
    
    def makeSendFunction(self):
        baseUrl = f'{self.svrName}_api_base_url'
        # 适配js的map解析，如果ret中有map，则将解析过后的object转换为map
        convert = ''
        for field in self.retStruct:
            if len(field[1]) > 4 and field[1][0:4] == 'map[':
                convert += f'   result.data.{field[0]} = new Map(Object.entries(result.data.{field[0]}));\n'
        if convert != '':
            convert = 'if (result.ok) {' + convert + '}\n'
        sendFunction = f'export function send{self.firstUpperCmd}(req: {self.firstUpperCmd}Req, callback: (ret: {{ data: {self.firstUpperCmd}Ret, ok: boolean }}) => void) {{\n'
        sendFunction += f'    UniPost({baseUrl} + \'{self.firstLowerCmd}\', req).then((res: UniResult) => {{\n'
        sendFunction += f'        const result: {{ data: {self.firstUpperCmd}Ret, ok: boolean }} = {{\n'
        sendFunction += f'            data: res.data as {self.firstUpperCmd}Ret,\n'
        sendFunction += f'            ok: res.ok\n'
        sendFunction += f'        }};\n'
        sendFunction += f'        {convert}\n'
        sendFunction += f'        callback(result);\n'
        sendFunction += f'    }});\n'
        sendFunction += f'}}\n'
        sendFunction += '\n'
        return sendFunction
        

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
with open(addr, "r", encoding="utf-8") as file:
    content = file.read()
    
# 读取上层文件夹名作为服务名
service_name = os.path.basename(os.path.dirname(addr))

# 匹配所有的msg
matches = re.findall(msg_pattern,content)

# 生成辅助结构体
helpers : list[Msg] = []
for match in matches:
    helpers.append(Msg(service_name,match[0],match[3],match[5]))
    
# 生成golang代码
golang_case = ""
golang_case += 'switch msg.Cmd(){\n'
for helper in helpers:
    golang_case += helper.makeCase()
golang_case += '}\n'
golang_on = ""
for helper in helpers:
    golang_on += helper.makeOnFunction()

# 生成typescript代码
typescript_head = 'import {UniPost, UniResult} from "../common/newSendHttp";'
typescript_interface = ""
for helper in helpers:
    typescript_interface += helper.makeTypeScriptInterface()
typescript_send = ""
typescript_send += f'const {service_name}_api_base_url = TODO;\n'
for helper in helpers:
    typescript_send += helper.makeSendFunction()
    
# 生成最新的结构体
latest_struct = ""
last = helpers[-1]
latest_struct += '===golang===\n'
latest_struct += last.makeCase()
latest_struct += last.makeOnFunction()
latest_struct += '===typescript===\n'
latest_struct += last.makeTypeScriptInterface()
latest_struct += last.makeSendFunction()
    

mode = input("""
1. 显示golang case代码
2. 复制golang case代码到剪贴板
3. 显示golang on代码
4. 复制golang on代码到剪贴板
5. 显示typescript interface代码
6. 复制typescript interface代码到剪贴板
7. 显示typescript send代码
8. 复制typescript send代码到剪贴板
9. 复制typescript全部代码到剪贴板
10.显示最新的结构体
请输入模式：__\b\b""")

import pyperclip

if mode == '1':
    print(golang_case)
elif mode == '2':
    pyperclip.copy(golang_case)
    print("已复制到剪贴板")
elif mode == '3':
    print(golang_on)
elif mode == '4':
    pyperclip.copy(golang_on)
    print("已复制到剪贴板")
elif mode == '5':
    print(typescript_interface)
elif mode == '6':
    pyperclip.copy(typescript_interface)
    print("已复制到剪贴板")
elif mode == '7':
    print(typescript_send)
elif mode == '8':
    pyperclip.copy(typescript_send)
    print("已复制到剪贴板")
elif mode == '9':
    pyperclip.copy(typescript_head + '\n' + typescript_interface + '\n' + typescript_send)
    print("已复制到剪贴板")
elif mode == '10':
    print(latest_struct)