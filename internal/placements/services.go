package placements

import (
	"reflect"
	"sync"
)

func partnerResponsesToMap(wg *sync.WaitGroup, response PartnerResponse, s map[uint][]PartnerResponseImp) {
	defer wg.Done()

	for _, val := range response.Imp {
		s[val.ID] = append(s[val.ID], val)
	}
}

func chanToSlice(ch interface{}) interface{} {
	chv := reflect.ValueOf(ch)
	slv := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(ch).Elem()), 0, 0)
	for {
		v, ok := chv.Recv()
		if !ok {
			return slv.Interface()
		}
		slv = reflect.Append(slv, v)
	}
}

func getMaxPrice(data []PartnerResponseImp) PartnerResponseImp {
	var resultIndex = 0
	maxVal := data[0].Price

	for i := range data {
		if data[i].Price > maxVal {
			maxVal = data[i].Price
			resultIndex = i
		}
	}

	return data[resultIndex]
}
