package poloniex

import (
	_"strconv"
	"testing"
	_"strings"

	"../../testhelpers"

	_"../../common/structs"
	"github.com/buger/jsonparser"
	_"github.com/sirupsen/logrus"
)

func TestParseSnapshotMessage(t *testing.T) {
	snapshotResponse, err := testhelpers.GetMsgFromFile("../../testhelpers/testdata/poloniex/test-channel-sub.json")
	if err != nil {
		t.Error("Could not read test file")
		return
	}

	// c := 0
	// d := 0
	// _, err = jsonparser.ArrayEach(snapshotResponse, func(value [] byte, datatype jsonparser.ValueType, offset int, err error) {
	// 	_, err = jsonparser.ArrayEach(value, func(value2 [] byte, datatype jsonparser.ValueType, offset int, err error) {
	// 		if (c != 2) {
	// 		t.Log(string(value2))
	// 		t.Log(datatype)
	// 		}

	// 		if (c == 2) {
	// 			t.Log(datatype)
	// 		_, err = jsonparser.ArrayEach(value2, func(value3 [] byte, datatype jsonparser.ValueType, offset int, err error) {
	// 			if (d != 1) {
	// 				t.Log(datatype)
	// 				//t.Log(string(value3))
	// 			}
	// 			d++
				
			
	// 			if (d == 1) {
	// 				t.Log(datatype)
	// 				/*
	// 				_, err = jsonparser.ArrayEach(value3 , func(value4 [] byte, datatype jsonparser.ValueType, offset int, err error) {
	// 					t.Log(string(value4))
						
	// 				})
	// 				*/
	// 			}
				
	// 		})
	// 		}
	// 		c++
	// 	})
	// })
	c := 0
	_, err = jsonparser.ArrayEach(snapshotResponse, func(value [] byte, datatype jsonparser.ValueType, offset int, err error) {
		_, err = jsonparser.ArrayEach(value, func(value2 [] byte, datatype jsonparser.ValueType, offset int, err error) {
			_, err = jsonparser.ArrayEach(value2, func(value3 [] byte, datatype jsonparser.ValueType, offset int, err error) {
				if (datatype == jsonparser.Object) {
					err = jsonparser.ObjectEach(value3, func(key [] byte, value4 [] byte, datatype jsonparser.ValueType, offset int) error {
						t.Log(string(key))
						return nil
					})
				c++
				}
			})
		})
	})

	//t.Log(snapshotResponse)
	return
}