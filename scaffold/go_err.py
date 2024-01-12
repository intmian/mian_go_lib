#!/usr/bin/env python
# -*- coding: utf-8 -*-
import os
import types
import re


class err_operator:
    def __init__(self, file_name: str, ori_code: str, ori_err: str) -> None:
        self.file_name = file_name
        self.ori_code = ori_code
        self.ori_err = ori_err

        # 将形如 has init 之类的错误信息转换为 ErrHasInit
        words = ori_err.split(" ")
        # 生成替换过后的错误常量名
        self.after_err = "Err"
        for word in words:
            # 如果首字母已经大写，就直接拼入
            if word[0].isupper():
                self.after_err += word
            else:
                self.after_err += word.capitalize()
        # 确认整体代换的代码
        if self.ori_code.startswith("errors.New"):
            self.replaced_code = self.after_err
        elif self.ori_code.startswith("errors.WithMessage"):
            self.replaced_code = self.ori_code.replace(
                '"' + self.ori_err + '"', self.after_err
            )
            self.replaced_code = self.replaced_code.replace(
                "errors.WithMessage", "errors.Join"
            )
        # 生成err文件中的行
        self.err_file_line = f'    {self.after_err} = ErrStr("{self.ori_err}")'


def find_err_and_pack(file: str) -> (map, str) or None:
    """
    找到所有包含
    fmt.Errorf("xxx",……)
    errors.New("xxx")
    errors.WithMessage(err, "xxx")
    如果第一个参数不为字符串则跳过
    如果fmt.Errorf，提取fmt.Errorf(".*", 和具体内容
    如果是errors.New，提取errors.New(".*") 和具体内容
    如果是errors.WithMessage(err, "xxx")，提取errors.WithMessage(err, "xxx") 和具体内容
    """
    # 文件名必须以.go结尾
    if not file.endswith(".go"):
        return None
    file_content = open(file, "r", encoding="utf-8").read()
    lines = file_content.split("\n")
    if len(lines) == 0:
        return None
    if "package" not in lines[0]:
        return None
    pack_name = lines[0].split("package ")[1]
    err_map = {}
    for line in lines:
        # if "fmt.Errorf" in line:
        #     err = re.search(r"fmt.Errorf\(\"(.*)\",", line)
        #     if err:
        #         e2 = err_operator(file, err[0], err[1])
        #         if not err_map.get(e2.ori_err):
        #             err_map[e2.after_err] = e2
        if "errors.New" in line:
            err = re.search(r"errors.New\(\"(.*)\"\)", line)
            if err:
                e2 = err_operator(file, err[0], err[1])
                if not err_map.get(e2.ori_err):
                    err_map[e2.after_err] = e2

        elif "errors.WithMessage" in line:
            err = re.search(r"errors\.WithMessage\(.*, *\"(.*)\"\)", line)
            if err:
                e2 = err_operator(file, err[0], err[1])
                if not err_map.get(e2.ori_err):
                    err_map[e2.after_err] = e2

    return err_map, pack_name


def main():
    err_file = "err_auto.go"
    pack_name = ""
    err_map = {}
    # os.chdir(os.path.dirname(os.path.realpath(__file__)))
    addr = "."
    # addr = "E:\\my_code_out\\platform\\mian_go_lib\\xnews"
    # all_after_err = []
    # 扫描当前目录所有非目录文件
    files = []
    for file in os.listdir(addr):
        if "." in file:
            # 完整路径
            file = os.path.join(addr, file)
            files.append(file)
    for file in files:
        errs = None
        results = find_err_and_pack(file)
        if results:
            errs, pack_name = results
            err_map[file] = errs
        # for err in errs:
        #     find = False
        #     for ae in all_after_err:
        #         if ae.after_err == err.after_err:
        #             find = True
        #             break
        #     if not find:
        #         all_after_err.append(err)

    # 汇总所有错误
    for file in err_map:
        print(f"file: {file}")
        with open(file, "r", encoding="utf-8") as f:
            strs = f.read()
            if '"github.com/pkg/errors"' in strs:
                print(f"import github.com/pkg/errors -> errors")
        err_map2 = err_map[file]
        for err in err_map2:
            value = err_map2[err]
            print(f"{value.ori_code} -> {value.replaced_code}, extract: {value.after_err}")

    choice = input("是否要生成错误包？(*/n)")
    if choice == "n":
        return

    # 修改文件
    for file in err_map:
        err_map2 = err_map[file]
        for err in err_map2:
            value = err_map2[err]
            file_content = open(file, "r", encoding="utf-8").read()
            # 此脚手架将会把 github.com/pkg/errors 替换为 errors，如果字符串中有可能会有问题，后续再看看
            file_content = file_content.replace('"github.com/pkg/errors"', '"errors"')
            file_content = file_content.replace(value.ori_code, value.replaced_code)
            open(file, "w", encoding="utf-8").write(file_content)

    err_file_head = "package " + pack_name + "\n\n"
    err_file_value_before = "const(\n"
    err_file_values = []
    if os.path.exists(err_file):
        with open(err_file, "r", encoding="utf-8") as f:
            lines = f.readlines()
            # 读取之前的错误
            for line in lines:
                if "= ErrStr" in line:
                    err_file_values.append(line)
    check_str = ""
    for err_file_value in err_file_values:
        check_str += err_file_value
    if 'ErrStr("nil")' not in check_str:
        err_file_values.append('    ErrNil = ErrStr("nil")\n')

    for file in err_map:
        err_map2 = err_map[file]
        for err in err_map2:
            value = err_map2[err]
            s = value.err_file_line  + "  // auto generated from " + value.file_name + "\n"
            if value.after_err not in check_str:
                err_file_values.append(s)

    err_file_value_after = ")\n\n"
    err_type = "type ErrStr string\n\n"
    err_func = "func (e ErrStr) Error() string { return string(e) }\n"

    with open(err_file, "w+", encoding="utf-8") as f:
        f.write(err_file_head)
        f.write(err_type)
        f.write(err_file_value_before)
        f.writelines(err_file_values)
        f.write(err_file_value_after)
        f.write(err_func)


if __name__ == "__main__":
    main()
