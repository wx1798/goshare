package goshare

import (
	// "log"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/mineralres/goshare/aproto"
)

func (p *Service) GetIndexMember(symbol *aproto.Symbol, retryCount int) ([]aproto.Symbol, error) {
	return getIndexMem(symbol)
	// var ret []aproto.Symbol;
	// return ret, nil
}

func getIndexMem(symbol *aproto.Symbol) ([]aproto.Symbol, error) {
	var ret []aproto.Symbol

	page_number := 1
	member_number := 0

	for true {
		address := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/corp/view/vII_NewestComponent.php?page=%d&indexid=%s", page_number, symbol.Code)
		// log.Println(address)
		page_number++
		doc, err := goquery.NewDocument(address)

		if err != nil {
			fmt.Println(err)
			return ret, err
		}

		b_empty := true

		doc.Find("#NewStockTable").Find("tbody").Find("tr").Each(func(i int, s *goquery.Selection) {
			if i > 0 {
				b_empty = false
				code := s.Find("div").Eq(0).Text()
				// fmt.Println(code)
				s, err := formatSymbol(code);
				if err == nil{
					ret = append(ret, s)
					member_number += 1
				}
			}
		})
		if b_empty == true || doc.Find("#page_form").Length() == 0 {
			break
		}
	}

	// log.Println(ret)
	return ret, nil
}

func formatSymbol(code string) (aproto.Symbol, error) {
	var ret aproto.Symbol
	if len(code) < 6{
		return ret, fmt.Errorf("error code %s", code)
	}

	switch code[0] {
	case '6':
		return aproto.Symbol{Exchange: aproto.ExchangeType_SSE, Code: code}, nil
	case '0':
		return aproto.Symbol{Exchange: aproto.ExchangeType_SZE, Code: code}, nil
	case '3':
		return aproto.Symbol{Exchange: aproto.ExchangeType_SZE, Code: code}, nil
	default:
		return ret, fmt.Errorf("error code %s", code)
	}
}
