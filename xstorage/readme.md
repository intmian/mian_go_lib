提供方便得通用存储，用于实现方便易用的通用存储、缓存，并能快速实现持久化存储。

为了通用起见，所有的存储都是基于key-value的形式，key为string，value为type interface{}。
建议将作为key的字符串以常量形式定义，以便于统一管理。

用于非高性能的存储，随时随地都能存储不需要单独建表，类似于需要访问的新闻关键词等，或者支持大系统的小模块肆意落盘等

类似于vscode，在配置json之上的操作、搜索、显示体系