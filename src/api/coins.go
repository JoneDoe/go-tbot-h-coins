package api

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

type positions struct {
	FIO   uint
	ID    uint
	Coins uint
}

// Account struct
type Account struct {
	Row, Col uint
	Balance  int
}

// GetCoins user coin balance
func GetCoins(userName string) (account Account) {
	sheet := getSheet()

	p := detectNeededDataPositions(sheet)

	for _, val := range sheet.Columns[p.ID] {
		if val.Value == userName {
			balance, _ := strconv.Atoi(sheet.Columns[p.Coins][val.Row].Value)
			account = Account{val.Row, p.Coins, balance}
			break
		}
	}

	return
}

// GetUserList coins recipients
func GetUserList() map[string]string {
	sheet := getSheet()

	p := detectNeededDataPositions(sheet)

	list := make(map[string]string)

	for _, val := range sheet.Columns[p.FIO] {
		if val.Row != 0 {
			list[sheet.Columns[p.ID][val.Row].Value] = val.Value
		}
	}

	return list
}

// ApplyTransaction apply move coins from-to accounts
func ApplyTransaction(accountSender Account, accountRecipient Account, amount int) {
	sheet := getSheet()

	log.Println(amount, accountSender)

	sheet.Update(int(accountSender.Row), int(accountSender.Col), strconv.Itoa(accountSender.Balance-amount))
	sheet.Update(int(accountRecipient.Row), int(accountRecipient.Col), strconv.Itoa(accountRecipient.Balance+amount))

	err := sheet.Synchronize()
	checkError(err)
}

func detectNeededDataPositions(sheet *spreadsheet.Sheet) *positions {
	p := new(positions)

	for _, val := range sheet.Rows[0] {
		switch val.Value {
		case "FIO":
			p.FIO = val.Column
		case "coins":
			p.Coins = val.Column
		case "ID":
			p.ID = val.Column
		}
	}

	return p
}

func getSheet() *spreadsheet.Sheet {
	data, err := ioutil.ReadFile(os.Getenv("GOOGLE_CLIENT_SECRET_PATH"))
	checkError(err)

	config, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	checkError(err)

	client := config.Client(context.TODO())
	service := spreadsheet.NewServiceWithClient(client)

	spreadsheet, err := service.FetchSpreadsheet(os.Getenv("HCOINS_SHEET_TOKEN"))
	checkError(err)
	sheet, err := spreadsheet.SheetByIndex(0)
	checkError(err)

	return sheet
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
