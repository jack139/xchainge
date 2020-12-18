package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	me      *user
	rootCmd = &cobra.Command{
		Use:   "xcli",
		Short: "xchainge client",
		Long:  "xcli is a client tool for xchainge",
	}
	dealCmd = &cobra.Command{	// 上链操作
		Use:   "deal",
		Short: "make a deal",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 5 {
				return errors.New("need more parameters")
			}
			action := args[0]
			assetsId := args[1]
			exchangeId := args[2]
			data := args[3]
			refer := args[4]
			return me.deal(action, assetsId, exchangeId, data, refer)
		},
	}
	authCmd = &cobra.Command{	// 上链操作
		Use:   "auth",
		Short: "query auth",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 5 {
				return errors.New("need more parameters")
			}
			action := args[0]
			assetsId := args[1]
			fromExchangeId := args[2]
			toExchangeId := args[3]
			refer := args[4]
			return me.auth(action, assetsId, fromExchangeId, toExchangeId, refer)
		},
	}

	generateKeyCmd = &cobra.Command{	// 生成密钥文件
		Use:   "generate",
		Short: "generate exchange key file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return me.generateKey()
		},
	}

	queryExchangeCmd = &cobra.Command{	// 查询 交易所 交易历史
		Use:   "queryExchange",
		Short: "query deals' history of exchange",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("need more parameters")
			}
			exchangeId := args[0]
			return me.query("exchange", exchangeId)
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
			return me.query("assets", assetsId)
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
			return me.query("refer", refer)
		},
	}

)

func init() {
	user, err := loadUserKeyFile()
	if err != nil {
		panic(err)
	}
	me = user

	rootCmd.AddCommand(dealCmd)
	rootCmd.AddCommand(authCmd)
	rootCmd.AddCommand(queryExchangeCmd)
	rootCmd.AddCommand(queryAssetsCmd)
	rootCmd.AddCommand(queryReferCmd)
	rootCmd.AddCommand(generateKeyCmd)	
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
