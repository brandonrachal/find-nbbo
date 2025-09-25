package exchange

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Data struct {
	Symbol   string
	Exchange string
	Bid      int
	Offer    int
}

type Exchange struct {
	dataMap map[string]map[string]*Data
}

func NewExchange() *Exchange {
	return &Exchange{
		dataMap: make(map[string]map[string]*Data),
	}
}

func (e *Exchange) Open(ctx context.Context, dataReader *bufio.Reader) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			readData, readErr := dataReader.ReadString('\n')
			if readErr != nil {
				if readErr == io.EOF {
					return nil
				}
				return readErr
			}
			if strings.HasPrefix(readData, "Q|") {
				prepareData := strings.TrimSuffix(readData, "\r\n")
				prepareData = strings.TrimSuffix(prepareData, "\n")
				e.addData(newData(prepareData))
			}
		}
	}
}

func (e *Exchange) addData(channelData *Data) {
	exchangeMap, ok := e.dataMap[channelData.Symbol]
	if !ok {
		e.dataMap[channelData.Symbol] = map[string]*Data{channelData.Exchange: channelData}
	} else {
		exchangeMap[channelData.Exchange] = channelData
	}
}

func (e *Exchange) NBBO(symbol string) (*string, error) {
	//Best bid = Max of (420.98, 420.93, 420.95) = 420.98
	//Best offer = Min of (421.00, 420.99, 421.05) = 420.99
	//NBBO = 420.98 @ 420.99
	bidData := e.MaxBid(symbol)
	offerData := e.MinOffer(symbol)
	if bidData != nil && offerData != nil {
		nbbo := fmt.Sprintf("%d@%d", bidData.Bid, offerData.Offer)
		return &nbbo, nil
	} else if bidData == nil {
		return nil, fmt.Errorf("couldn't find a bid for symbol %s", symbol)
	} else {
		return nil, fmt.Errorf("couldn't find a offer for symbol %s", symbol)
	}
}

func (e *Exchange) MaxBid(symbol string) *Data {
	var bestBid *Data
	exchangeMap, ok := e.dataMap[symbol]
	if !ok {
		return nil
	}
	for _, exchangeData := range exchangeMap {
		if bestBid == nil {
			bestBid = exchangeData
		} else if exchangeData.Bid >= bestBid.Bid {
			bestBid = exchangeData
		}
	}
	return bestBid
}

func (e *Exchange) MinOffer(symbol string) *Data {
	var bestOffer *Data
	exchangeMap, ok := e.dataMap[symbol]
	if !ok {
		return nil
	}
	for _, exchangeData := range exchangeMap {
		if bestOffer == nil {
			bestOffer = exchangeData
		} else if exchangeData.Offer <= bestOffer.Offer {
			bestOffer = exchangeData
		}
	}
	return bestOffer
}

func (e *Exchange) Debug() {
	jsonBytes, jsonBytesErr := json.MarshalIndent(e.dataMap, "", "  ")
	if jsonBytesErr != nil {
		fmt.Println(jsonBytesErr)
		return
	}
	fmt.Println("Exchange Data:")
	fmt.Println(string(jsonBytes))
}

func newData(data string) *Data {
	dataList := strings.Split(data, "|")
	return &Data{
		Symbol:   dataList[1],
		Exchange: dataList[2],
		Bid:      toInt(dataList[3]),
		Offer:    toInt(dataList[4]),
	}
}

func toInt(data string) int {
	integer, _ := strconv.Atoi(data)
	return integer
}
