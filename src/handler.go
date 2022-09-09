package app

import (
	"errors"
	"fmt"
	"log"

	api "go-tbot-h-coins/src/api"
	blockchain "go-tbot-h-coins/src/blockchain"
)

type transStruct struct {
	Date   string `json:"date"`
	From   string `json:"from"`
	To     string `json:"to"`
	Amount int    `json:"amount"`
}

type handler struct {
	bc *blockchain.Blockchain
}

// Handler handler
var Handler handler

func (h *handler) CoinsTransfer(transaction CoinsTransaction) (err error) {
	accountSender := api.GetCoins(User.Get(request.ID).PhoneNumber)

	if transaction.Amount > accountSender.Balance {
		err := errors.New("Сумма превышает баланс")
		return err
	}

	accountRecipient := api.GetCoins(transaction.Recipient)

	api.ApplyTransaction(accountSender, accountRecipient, transaction.Amount)

	h.rectransaction(transaction)

	/* bc := blockchain.NewBlockchain()
	bc.AddBlock("ddd send 22 BTC to 3333")

	list := api.GetUserList()
	rec := list[transaction.Recipient]
	log.Println(fmt.Sprintf("%s send %d BTC to %s", transaction.Sender, transaction.Amount, rec)) */

	/* fp := path.Join("log", "trans.json")
	fp. */

	/* filename := path.Join("log", "trans.txt")

	//file, _ := os.Create(filename)
	file, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer file.Close()

	jsonBlob := []byte(`{"Date":"Dummy","From":"Vova","To":"Valja", "Amount":150}`)
	rankings := transStruct{}
	err = json.Unmarshal(jsonBlob, &rankings)
	if err != nil {
		// nozzle.printError("opening config file", err.Error())
	}

	rankingsJson, _ := json.Marshal(rankings)
	file.Write(rankingsJson)
	//err = ioutil.WriteFile("output.json", rankingsJson, 0644)
	fmt.Printf("%+v", rankings) */

	transaction = CoinsTransaction{0, "", ""}

	return nil
}

func (h *handler) rectransaction(transaction CoinsTransaction) {
	list := api.GetUserList()
	rec := list[transaction.Recipient]
	sender := transaction.Sender

	//fmt.Printf("%+v", transaction)

	log.Println(fmt.Sprintf("%s send %d BTC to %s", sender, transaction.Amount, rec))

	h.bc.AddBlock(fmt.Sprintf("%s send %d BTC to %s", sender, transaction.Amount, rec))
	//defer bc.db.Close()
}
