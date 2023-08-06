package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var r, t, d string

type FileInfo struct {
	FullName string
	Ext      string
	NewName  string
}

func main() {

	// cli root cmd
	var rootCmd = &cobra.Command{
		Use:   "extChange -R regex -T new_ext [-D dir]",
		Short: "Change file ext matched REGEX to NEW_EXT, if DIR is empty, use current directory.",
		Run:   Run,
	}
	rootCmd.Flags().StringVarP(&r, "regex", "R", "", "matched EXT regex")
	rootCmd.Flags().StringVarP(&t, "ext", "T", "", "new ext")
	rootCmd.Flags().StringVarP(&d, "dir", "D", "", "target dir")
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal("run error", err)
		os.Exit(1)
	}
}

func Run(cmd *cobra.Command, args []string) {
	if r == "" {
		log.Fatal("需要指定需要修改的ext")
		return
	}
	re, err := regexp.Compile(r)
	if err != nil {
		log.Fatalf("regex有误:%s", r)
		return
	}
	if t == "" {
		log.Fatal("需要指定目标ext")
		return
	}
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal("无法获取当前目录：", err)
		return
	}
	if len(d) > 1 && d[1] == ':' || path.IsAbs(d) {
	} else {
		d = path.Join(dir, d)
	}
	reader := bufio.NewReader(os.Stdin)
	mf := matchedFiles(re)
	if len(mf) == 0 {
		log.Println("目录为空或无匹配文件")
		return
	}
	fmt.Printf("即将执行目录 %s 下文件的重命名，你确定要执行操作吗？(Y/N):\n", d)
	for i, f := range mf {
		fmt.Printf("%d: %s -> %s\n", i, f.FullName, f.NewName)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if strings.EqualFold(input, "Y") || strings.EqualFold(input, "Yes") {
		// 遍历所有文件
		for _, f := range mf {
			// 重命名文件
			err = os.Rename(path.Join(d, f.FullName), path.Join(d, f.NewName))
			if err != nil {
				log.Fatalf("无法重命名文件 %s：%s\n", f.FullName, err)
			} else {
				log.Printf("已重命名文件 %s 为 %s\n", f.FullName, f.NewName)
			}
		}
	} else {
		log.Println("User Cancel, Exit")
		return
	}

}

func matchedFiles(re *regexp.Regexp) []*FileInfo {
	// 读取目录下的所有文件和文件夹
	files, err := os.ReadDir(d)
	if err != nil {
		log.Fatal("无法读取当前目录：", err)
		return nil
	}
	mf := make([]*FileInfo, 0)
	for _, file := range files {
		if !file.IsDir() {
			// 获取文件名和扩展名
			fullName := file.Name()
			ext := filepath.Ext(fullName)[1:]
			newName := strings.TrimSuffix(fullName, ext) + t
			info := &FileInfo{
				FullName: file.Name(),
				Ext:      filepath.Ext(file.Name())[1:],
				NewName:  newName,
			}
			// 检查扩展名
			if matched := re.Match([]byte(info.Ext)); matched && ext != t {
				mf = append(mf, info)
			}
		}
	}
	return mf
}
