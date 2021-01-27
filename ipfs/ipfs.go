package ipfs


import (
	"fmt"
	"strings"
	"io"

	shell "github.com/ipfs/go-ipfs-api"
)


const HOST = "localhost:5001"


func Add(filedata []byte) (string, error) {
	// 连接api
	sh := shell.NewShell("localhost:5001")

	// 添加内容
	cid, err := sh.Add(strings.NewReader(string(filedata)))
	if err != nil {
		return "", fmt.Errorf("IPFS error: %s\n", err)
	}

	return cid, nil
}

func Get(cid string) ([]byte, error) {
	// 连接api
	sh := shell.NewShell("localhost:5001")

	// 获取文件内容
	data, err := sh.Cat(cid)
	if err != nil {
		return nil, fmt.Errorf("IPFS error: %s\n", err)
	}
	defer data.Close()

	// 使用缓存读出文件
	var dataBuf []byte
	longBuf := make([]byte, 1024*20)

	for {
		sz, err := data.Read(longBuf)
		if err != nil {
			if err == io.EOF {
				if sz>0 { // EOF 此时有可能还读出了数据
					fmt.Printf("EOF: n = %d\n", sz)
					dataBuf = append(dataBuf, longBuf[:sz]...)
				}
				break
			}
			return nil, fmt.Errorf("IPFS error: %s\n", err)
		}
		//fmt.Printf("%d %s\n", sz, longBuf)
		dataBuf = append(dataBuf, longBuf[:sz]...)
	}

	return dataBuf, nil
}