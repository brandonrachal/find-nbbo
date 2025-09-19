package exchange

import (
	"bufio"
	"context"
	_ "embed"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

//go:embed testdata/data.txt
var fileData string

func TestExchange(t *testing.T) {
	req := require.New(t)

	dataReader := bufio.NewReader(strings.NewReader(fileData))
	ctx := context.Background()
	exchange := NewExchange()
	go func() {
		err := exchange.Open(ctx, dataReader)
		if err != nil {
			req.NoError(err)
		}
	}()

	// Ensure data is processed
	time.Sleep(1 * time.Second)
	//exchange.Debug()

	symbol := "DELL"
	expectedNBBO := "1072@1090"
	nBBO, nBBOErr := exchange.NBBO(symbol)
	req.NoError(nBBOErr)
	req.Equal(expectedNBBO, *nBBO)

	symbol = "IBM"
	expectedNBBO = "1090@1089"
	nBBO, nBBOErr = exchange.NBBO(symbol)
	req.NoError(nBBOErr)
	req.Equal(expectedNBBO, *nBBO)
}
