package main

/*
	主程序：

	编译：
	go build main.go

	运行：
	./xcli
*/

import (
	"xchainge/client"
	"xchainge/ipfs"
	"xchainge/http"

	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	me      *client.User
	rootCmd = &cobra.Command{
		Use:   "xcli",
		Short: "xchainge client",
		Long:  "xcli is a client tool for xchainge",
	}
	initCmd = &cobra.Command{ // 生成key
		Use:   "init",
		Short: "make user key",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			// 如果 path不存在，会创建密钥文件
			_, err := client.GenUserKey(path)
			if err!=nil {
				return err
			}
			return nil
		},
	}
	dealCmd = &cobra.Command{	// 交易上链操作
		Use:   "deal",
		Short: "make a deal",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 4 {
				return errors.New("need more parameters")
			}
			action := args[0]
			assetsId := args[1]
			data := args[2]
			refer := args[3]
			respBytes, err := me.Deal(action, assetsId, data, refer)
			if err==nil {
				fmt.Printf("Deal ==> %s\n", respBytes)
			}
			return err
		},
	}
	authRequestCmd = &cobra.Command{	// 上链操作，请求授权
		Use:   "authReq",
		Short: "Request authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			fromExchangeId := args[0]
			dealId := args[1] // 请求授权的 dealID
			respBytes, err := me.AuthRequest(fromExchangeId, dealId)
			if err==nil {
				fmt.Printf("AuthReq ==> %s\n", respBytes)
			}
			return err
		},
	}
	authResponseCmd = &cobra.Command{	// 上链操作，响应授权
		Use:   "authResp",
		Short: "Respond to authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			authId := args[0] // 响应授权的 authID
			respBytes, err := me.AuthResponse(authId)
			if err==nil {
				fmt.Printf("AuthResp ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryDealCmd = &cobra.Command{	// 查询 交易所 交易历史
		Use:   "queryDeal",
		Short: "query deals' history of exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			respBytes, err := me.Query("deal", "_")
			if err==nil {
				fmt.Printf("Deal ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryAuthCmd = &cobra.Command{	// 查询 请求授权 历史
		Use:   "queryAuth",
		Short: "query requests of authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			respBytes, err := me.Query("auth", "_")
			if err==nil {
				fmt.Printf("Auth ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryAssetsCmd = &cobra.Command{	// 查询 资产 交易历史
		Use:   "queryAssets",
		Short: "query deals' history of Assets",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			assetsId := args[0]
			respBytes, err := me.Query("assets", assetsId)
			if err==nil {
				fmt.Printf("Assets ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryReferCmd = &cobra.Command{	// 查询 Refer 交易历史
		Use:   "queryRefer",
		Short: "query deals' history of Refer",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			refer := args[0]
			respBytes, err := me.Query("refer", refer)
			if err==nil {
				fmt.Printf("Refer ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryTxCmd = &cobra.Command{	// 查询 指定交易
		Use:   "queryTx",
		Short: "query deals by DealID",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			exchangeId := args[0]
			realId := args[1]
			respBytes, err := me.QueryTx(exchangeId, realId)
			if err==nil {
				fmt.Printf("Tx ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryRawCmd = &cobra.Command{	// 查询 指定raw block
		Use:   "queryRaw",
		Short: "query block raw data by DealID",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			exchangeId := args[0]
			realId := args[1]
			respBytes, err := me.QueryRawBlock(exchangeId, realId)
			if err==nil {
				fmt.Printf("Raw ==> %s\n", respBytes)
			}
			return err
		},
	}

	creditCmd = &cobra.Command{	// CR上链操作
		Use:   "credit",
		Short: "make credits",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			if len(args) < 3 {
				return errors.New("need more parameters")
			}
			action := args[0]
			data := args[1]
			num := args[2]
			respBytes, err := me.Credit(action, data, num)
			if err==nil {
				fmt.Printf("Credit ==> %s\n", respBytes)
			}
			return err
		},
	}

	queryCreditCmd = &cobra.Command{	// 查询 CR
		Use:   "queryCR",
		Short: "query CR history",
		RunE: func(cmd *cobra.Command, args []string) error {
			path, _ := cmd.Flags().GetString("home")
			me, err := client.GetMe(path)
			if err!=nil {
				return err
			}
			respBytes, err := me.Query("credit", "_")
			if err==nil {
				fmt.Printf("Credit ==> %s\n", respBytes)
			}
			return err
		},
	}

	httpCmd = &cobra.Command{	// 启动http服务
		Use:   "http <port> <user_path>",
		Short: "start http service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("need port number")
			} 
			port := args[0]
			userPath := args[1]
			http.RunServer(port, userPath)
			// 不会返回
			return nil
		},
	}

	ipfsAdd = &cobra.Command{	// ipfs add
		Use:   "ipfsAdd <text>",
		Short: "add text to ipfs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need text")
			} 
			text := args[0]
			cid, err := ipfs.Add([]byte(text))
			if err==nil {
				fmt.Printf("IPFS add ==> %s\n", cid)
			}
			return err
		},
	}

	ipfsGet = &cobra.Command{	// ipfs get
		Use:   "ipfsGet <cid>",
		Short: "get text from IPFS",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("need cid")
			} 
			cid := args[0]
			respBytes, err := ipfs.Get(cid)
			if err==nil {
				fmt.Printf("IPFS get ==> %s\n", respBytes)
			}
			return err
		},
	}

)

func init() {
	initCmd.Flags().StringP("home", "", "", "密钥文件路径");
	dealCmd.Flags().StringP("home", "", "", "密钥文件路径");
	authRequestCmd.Flags().StringP("home", "", "", "密钥文件路径")
	authResponseCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryDealCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryAssetsCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryReferCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryAuthCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryTxCmd.Flags().StringP("home", "", "", "密钥文件路径")
	queryRawCmd.Flags().StringP("home", "", "", "密钥文件路径")
	creditCmd.Flags().StringP("home", "", "", "密钥文件路径");
	queryCreditCmd.Flags().StringP("home", "", "", "密钥文件路径")

	initCmd.MarkFlagRequired("home")
	dealCmd.MarkFlagRequired("home")
	authRequestCmd.MarkFlagRequired("home")
	authResponseCmd.MarkFlagRequired("home")
	queryDealCmd.MarkFlagRequired("home")
	queryAssetsCmd.MarkFlagRequired("home")
	queryReferCmd.MarkFlagRequired("home")
	queryAuthCmd.MarkFlagRequired("home")
	queryTxCmd.MarkFlagRequired("home")
	queryRawCmd.MarkFlagRequired("home")
	creditCmd.MarkFlagRequired("home")
	queryCreditCmd.MarkFlagRequired("home")

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(dealCmd)
	rootCmd.AddCommand(authRequestCmd)
	rootCmd.AddCommand(authResponseCmd)
	rootCmd.AddCommand(queryDealCmd)
	rootCmd.AddCommand(queryAssetsCmd)
	rootCmd.AddCommand(queryReferCmd)
	rootCmd.AddCommand(queryAuthCmd)
	rootCmd.AddCommand(queryTxCmd)
	rootCmd.AddCommand(queryRawCmd)
	rootCmd.AddCommand(httpCmd)
	rootCmd.AddCommand(creditCmd)
	rootCmd.AddCommand(queryCreditCmd)
	rootCmd.AddCommand(ipfsAdd)
	rootCmd.AddCommand(ipfsGet)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
