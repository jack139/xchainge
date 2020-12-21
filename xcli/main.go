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
	dealCmd = &cobra.Command{	// 上链操作
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
	authCmd = &cobra.Command{	// 上链操作
		Use:   "auth",
		Short: "query auth",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 4 {
				return errors.New("need more parameters")
			}
			action := args[0]
			assetsId := args[1]
			toExchangeId := args[2]
			refer := args[3]
			return me.Auth(action, assetsId, toExchangeId, refer)
		},
	}

	queryExchangeCmd = &cobra.Command{	// 查询 交易所 交易历史
		Use:   "queryExchange",
		Short: "query deals' history of exchange, '_' for local exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			exchangeId := args[0]
			return me.Query("exchange", exchangeId)
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
			return me.Query("assets", assetsId)
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
			return me.Query("refer", refer)
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
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(queryExchangeCmd)
	rootCmd.AddCommand(queryAssetsCmd)
	rootCmd.AddCommand(queryReferCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
