package internal

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

var (
	// 从字符串 "#INCLUDE [dir1/dir2/file.txt] @SECTION_1" 中提取出 "dir1/dir2/file.txt" 和 "@SECTION_1"
	// 合并指定文件中的指定节点的内容到当前文件
	regexInclude = regexp.MustCompile(`#INCLUDE\s*\[([^\n]+)\]\s*(@[^\n]+)`)

	// 从字符串 "#INSERT [dir1/dir2/file.txt] @SECTION_1" 中提取出 "dir1/dir2/file.txt"
	// 合并指定文件中的所有内容到当前文件
	regexInsert = regexp.MustCompile(`#INSERT\s*\[([^\n]+)\]\s*$`)

	// 从字符串 "[@MAIN]" 中提取出 "@MAIN"
	regexFunction = regexp.MustCompile(`^\[([^\n]+)\]\s*$`)
)

func LoadFile(file string) (*Script, error) {
	var r, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	return Load(r)
}

func Load(r io.Reader) (*Script, error) {
	var lines, err = ExpandScript(Read(r))
	if err != nil {
		return nil, err
	}

	var f *Function
	var script = NewScript()

	for _, line := range lines {

		fmt.Println(line)

		if line[0] == '[' {
			var match = regexFunction.FindStringSubmatch(line)
			if len(match) > 0 {
				if f != nil {
					script.Add(f)
				}

				f = NewFunction(match[1])
				continue
			}
		}

		if f != nil {
			f.Add(line)
		}
	}

	if f != nil {
		script.Add(f)
	}

	return script, nil
}

// ExpandScript 处理 #INSERT 语句和 #INCLUDE 语句
func ExpandScript(lines []string) ([]string, error) {
	var nLines []string
	for _, line := range lines {
		if SkipLine(line) {
			continue
		}

		if line[0] == KeyPrefix {
			if strings.HasPrefix(line, KeyInsert) {
				var match = regexInsert.FindStringSubmatch(line)
				var insertLines, err = ReadFile(match[1])
				if err != nil {
					return nil, err
				}
				insertLines, err = ExpandScript(insertLines)
				if err != nil {
					return nil, err
				}
				nLines = append(nLines, insertLines...)
				continue
			} else if strings.HasPrefix(line, KeyInclude) {
				var match = regexInclude.FindStringSubmatch(line)
				var insertLines, err = Include(match[1], match[2])
				if err != nil {
					return nil, err
				}
				nLines = append(nLines, insertLines...)
				continue
			}
		}
		nLines = append(nLines, line)
	}
	return nLines, nil
}

func ReadFile(file string) ([]string, error) {
	var r, err = os.Open(file)
	if err != nil {
		return nil, err
	}
	return Read(r), nil
}

func Read(r io.Reader) []string {
	var lines []string

	var scanner = bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, string(RemoveBOM(scanner.Bytes())))
	}
	return lines
}

// Include 读取指定文件中的指定片断
func Include(file, name string) ([]string, error) {
	var lines, err = ReadFile(file)
	if err != nil {
		return nil, err
	}

	name = "[" + name + "]"

	var stat = 0

	var nLines []string
	for _, line := range lines {
		if SkipLine(line) {
			continue
		}
		switch stat {
		case 0:
			if line[0] == '[' && strings.HasPrefix(line, name) {
				stat = 1
			}
		case 1:
			if line[0] == '{' {
				stat = 2
			}
		case 2:
			if line[0] == '}' {
				return nLines, nil
			}

			nLines = append(nLines, line)
		}
	}
	return nil, errors.New("syntax error:" + file)
}
func SkipLine(line string) bool {
	if line == "" {
		return true
	}
	if line[0] == KeyComment {
		return true
	}
	return false
}
