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
	dealCmd = &cobra.Command{	// 交易上链操作
		Use:   "deal",
		Short: "make a deal",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return errors.New("need more parameters")
			}
			action := args[0]
			assetsId := args[1]
			data := args[2]
			refer := args[3]
			return me.Deal(action, assetsId, data, refer)
		},
	}
	authRequestCmd = &cobra.Command{	// 上链操作，请求授权
		Use:   "authReq",
		Short: "Request authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("need more parameters")
			}
			fromExchangeId := args[0]
			dealId := args[1] // 请求授权的 dealID
			return me.AuthRequest(fromExchangeId, dealId)
		},
	}
	authResponseCmd = &cobra.Command{	// 上链操作，响应授权
		Use:   "authResp",
		Short: "Respond to authorization",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			authId := args[0] // 响应授权的 authID
			return me.AuthResponse(authId)
		},
	}

	queryDealCmd = &cobra.Command{	// 查询 交易所 交易历史
		Use:   "queryDeal",
		Short: "query deals' history of exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
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

	httpCmd = &cobra.Command{	// 启动http服务
		Use:   "http",
		Short: "start http service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("need port number")
			} 
			http.RunServer(args[0])
			// 不会返回
			return nil
		},
	}

)

func init() {
	user, err := client.LoadOrGenUserKey()
	if err != nil {
		panic(err)
	}
	me = user

	rootCmd.AddCommand(dealCmd)
	rootCmd.AddCommand(authRequestCmd)
	rootCmd.AddCommand(authResponseCmd)
	rootCmd.AddCommand(queryDealCmd)
	rootCmd.AddCommand(queryAssetsCmd)
	rootCmd.AddCommand(queryReferCmd)
	rootCmd.AddCommand(queryAuthCmd)
	rootCmd.AddCommand(queryTxCmd)
	rootCmd.AddCommand(httpCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
