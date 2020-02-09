package coinbasepro

import(
    "testing"
    "io/ioutil"
    "encoding/json"

    "github.com/buger/jsonparser"
)

func getMsgFromFile(filepath string) ([]byte, error) {
    jsonFile, err := ioutil.ReadFile(filepath)
    if err != nil {
        return nil, err
    }

    return jsonFile, nil
}

func TestParseL2Message(t *testing.T) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-snapshot-response.json")
    if err != nil {
        t.Error("Could not read test file.")
    }

    var l2SnapshotMessage L2SnapshotMessage

    err = json.Unmarshal(jsonFile, &l2SnapshotMessage)
    if err != nil {
        t.Error("Unmarshal error.")
    }

    if l2SnapshotMessage.MessageType != "snapshot" {
        t.Errorf("Expected MessageType: %s, got %s", "snapshot", l2SnapshotMessage.MessageType)
    }
    if l2SnapshotMessage.ProductId != "ETH-USD" {
        t.Errorf("Expected ProductId: %s, got %s", "ETH-USD", l2SnapshotMessage.ProductId)
    }

    if l2SnapshotMessage.Asks[0][0] != "188.52" && l2SnapshotMessage.Asks[0][1] != "14.06613634" {
        t.Errorf("First Ask Price mismatched. Got %s, %s", l2SnapshotMessage.Asks[0][0], l2SnapshotMessage.Asks[0][1])
    }
    if l2SnapshotMessage.Bids[0][0] != "188.51" && l2SnapshotMessage.Bids[0][1] != "0.26060000" {
        t.Errorf("First Bid Price mismatched. Got %s, %s", l2SnapshotMessage.Bids[0][0], l2SnapshotMessage.Bids[0][1])
    }
}

func BenchmarkL2MessageJSON(b *testing.B) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-snapshot-response.json")
    if err != nil {
        b.Error("Could not read test file.")
    }
    b.ResetTimer()

    var l2SnapshotMessage L2SnapshotMessage

    for i := 0; i < b.N; i++ {
        err = json.Unmarshal(jsonFile, &l2SnapshotMessage)
        if err != nil {
            b.Error("Unmarshal error.")
        }
        b.Log(l2SnapshotMessage.Asks[0][0], ",", l2SnapshotMessage.Asks[0][1])
        b.Log(l2SnapshotMessage.Asks[1][0], ",", l2SnapshotMessage.Asks[1][1])
        b.Log(l2SnapshotMessage.Asks[2][0], ",", l2SnapshotMessage.Asks[2][1])
        
        b.Log(l2SnapshotMessage.Bids[0][0], ",", l2SnapshotMessage.Bids[0][1])
        b.Log(l2SnapshotMessage.Bids[1][0], ",", l2SnapshotMessage.Bids[1][1])
        b.Log(l2SnapshotMessage.Bids[2][0], ",", l2SnapshotMessage.Bids[2][1])
    }
}

func BenchmarkL2MessageJSONParser(b *testing.B) {
    jsonFile, err := getMsgFromFile("../../testdata/coinbasepro/test-l2-snapshot-response.json")
    if err != nil {
        b.Error("Could not read test file.")
    }
    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        v1, _:= jsonparser.GetString(jsonFile, "asks", "[0]", "[0]")
        v2, _:= jsonparser.GetString(jsonFile, "asks", "[0]", "[1]")
        b.Log(v1, ",", v2)

        v11, _:= jsonparser.GetString(jsonFile, "asks", "[1]", "[0]")
        v22, _:= jsonparser.GetString(jsonFile, "asks", "[1]", "[1]")
        b.Log(v11, ",", v22)

        v111, _:= jsonparser.GetString(jsonFile, "asks", "[2]", "[0]")
        v222, _:= jsonparser.GetString(jsonFile, "asks", "[2]", "[1]")
        b.Log(v111, ",", v222)

        v3, _ := jsonparser.GetString(jsonFile, "bids", "[0]", "[0]")
        v4, _ := jsonparser.GetString(jsonFile, "bids", "[0]", "[1]")
        b.Log(v3, ",", v4)

        v33, _ := jsonparser.GetString(jsonFile, "bids", "[1]", "[0]")
        v44, _ := jsonparser.GetString(jsonFile, "bids", "[1]", "[1]")
        b.Log(v33, ",", v44)

        v333, _ := jsonparser.GetString(jsonFile, "bids", "[2]", "[0]")
        v444, _ := jsonparser.GetString(jsonFile, "bids", "[2]", "[1]")
        b.Log(v333, ",", v444)
    }

}