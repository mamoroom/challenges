package main

import (
	"os"
	"log"
	"bytes"
	"sort"
	"strconv"
	"encoding/csv"
	"fmt"
	"time"

	"go-tamboon/cipher"	
	"go-tamboon/payment"	
)

const (
	currency = "thb"
)

type envelope struct {
	name string
	amount int64
	ccNumber string
	cvv string
	expMonth time.Month
	expYear int
}

type donor struct {
	name string
	amount int64
}


func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	omisePayment, e := payment.NewOmisePayment(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"))
	if e != nil {
		log.Fatal(e)
	}

	reader, e := readCsv(`data/fng.1000.csv.rot128`)
	if e != nil {
		log.Fatal(e)
	}

	lines, e := reader.ReadAll()
	if e != nil {
		log.Fatal(e)
	}
			
	var totalReceived int64 = 0
	var successfullyDonated int64 = 0
	var faultyDonation int64 = 0

	var donors []donor

	fmt.Print("performing donations")
	for i, record := range lines {
		fmt.Print(".")
		if i == 0 {
			// skip header line
			continue
		}

		// if i % 20 == 0 {
		// 	time.Sleep(time.Second * 3)
		// 	fmt.Println(".")
		// }

        // record, err := reader.Read()
        // if err == io.EOF {
        //     break
        // } else if err != nil {
        //     panic(err)
        // }

        envelope, e := parseEnvelope(record)
        if e != nil {
        	log.Fatal(e)
        }

        totalReceived += envelope.amount
        amount, e := execPayment(envelope, omisePayment)
        if e != nil {
        	// invalid card infomation
        	faultyDonation += envelope.amount
        	continue
        }

        successfullyDonated += amount
        donors = append(donors, donor{
        	name: envelope.name,	
        	amount: amount,
        })
	}
    var averagePerPersion float64 = float64(successfullyDonated) / float64(len(donors))
    fmt.Print("\ndone.\n\n")
    fmt.Printf("(currency is: %s)\n", currency)
    fmt.Printf("[total received] %d\n[successfully donated] %d\n[faulty donation] %d\n[average per persion] %.2f\n", totalReceived, successfullyDonated, faultyDonation, averagePerPersion)

    sort.Slice(donors, func(i, j int) bool {
    	return donors[i].amount > donors[j].amount
    })
    fmt.Printf("\n[top donors]\n")
    for i, donor := range donors {
    	if i > 2 {
    		break
    	}
    	fmt.Printf("%s\n", donor.name)
    }
}


func execPayment(envelope *envelope, p *payment.OmisePayment) (int64, error) {
    token, e := p.CreateToken(envelope.name, envelope.ccNumber, envelope.expMonth, envelope.expYear, envelope.cvv)
  	if e != nil {
  		return 0, e
	}

    charge, e := p.Charge(envelope.amount, currency, token)
	if e != nil {
  		return 0, e
	}
    //log.Printf("charge: %s  amount: %s %d\n", charge.ID, charge.Currency, charge.Amount)

    return charge.Amount, nil
}

func readCsv(path string) (*csv.Reader, error){ 
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fi.Size()
	// log.Println("Read file size is %d", fileSize)

	rot128Reader, err := cipher.NewRot128Reader(file)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte,fileSize,fileSize)
	_, err = rot128Reader.Read(buffer)
	if err != nil {
		return nil, err
	}

	return csv.NewReader(bytes.NewReader(buffer)), nil
}

func parseEnvelope(record []string) (*envelope, error) {
    amount, err := strconv.ParseInt(record[1], 10, 64)
    if err != nil {
    	return nil, err
    }
    month, err := strconv.Atoi(record[4])
    if err != nil {
    	return nil, err
    }
    year, err := strconv.Atoi(record[5])
    if err != nil {
    	return nil, err
    }

    return &envelope{
		name: record[0],
		amount: amount,
		ccNumber: record[2],
		cvv: record[3],
		expMonth: time.Month(month),
		expYear: year,
    }, nil
}

