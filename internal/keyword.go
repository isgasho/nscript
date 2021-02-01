package internal

// 默认函数
const (
	// 退出脚本
	FuncExit = "@EXIT"

	// 退出脚本
	FuncClose = "@CLOSE"
)

// 默认指令
const (
	CmdGoto  = "GOTO"
	CmdBreak = "BREAK"
)

const (
	KeyPrefix = '#'
)

// 关键字
const (
	// 注释
	KeyComment = ';'

	// IF 语句，执行判断
	KeyIf = "#IF"

	// SAY 语句，输出内容
	KeySay = "#SAY"

	// ELSESAY 语句，输出内容
	KeyElseSay = "#ELSESAY"

	// ACT 语句，执行操作
	KeyAct = "#ACT"

	// ELSEACT 语句，执行操作
	KeyElseAct = "#ELSEACT"

	// 将指定脚本文件中的所有内容引入到当前脚本中
	// 示例：
	// #INSERT [dir1/dir2/file.txt]
	KeyInsert = "#INSERT"

	// 将指定脚本文件中的特定脚本片断(函数)的内容引入到当前脚本中，不包含片断(函数)名
	// 示例：
	// #INCLUDE [dir1/dir2/file.txt] @SECTION_1
	KeyInclude = "#INCLUDE"
)
