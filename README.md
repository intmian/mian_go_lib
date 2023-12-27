# mian_go_lib

自用的go小轮子

| 子包         | 描述                         | 包含                                             |
|------------|----------------------------|------------------------------------------------|
| cipher     | 密码学包                       | aes、sha256                                     |
| cmd_server | tcp服务器包                    |                                                |
| menu       | 通用菜单包                      |                                                |
| misc       | 杂项包                        | 一堆小工具                                          |
| spider     | 爬虫包                        | 百度新闻、大盘、彩票                                     |
| xpush      | 推送包                        | 邮件、pushdeer、钉钉 sdk                             |
| xlog       | 通用日志包                      | 可以和push结合使用，并自主替换命令行输出                         |
| xstorage   | 线程安全、使用方便得，支持复杂配置的自落盘cache | 支持int、float、bool、string和对应的slice，包含一个简单的web外包装 |

可以参考test中的集成测试使用
