package main

import (
	"encoding/json"
	"fmt"
	"github.com/EdgeCast/vflow/ipfix"
	"github.com/gobeam/stringy"
	"strconv"
)

var ipfixChMap = map[ipfix.FieldType]string{
	ipfix.Uint8:                "UInt8",
	ipfix.Uint16:               "UInt16",
	ipfix.Uint32:               "UInt32",
	ipfix.Uint64:               "UInt64",
	ipfix.Int8:                 "Int8",
	ipfix.Int16:                "Int16",
	ipfix.Int32:                "Int32",
	ipfix.Int64:                "Int64",
	ipfix.Float32:              "Float32",
	ipfix.Float64:              "Float64",
	ipfix.Boolean:              "Bool",
	ipfix.MacAddress:           "String",
	ipfix.OctetArray:           "String",
	ipfix.String:               "String",
	ipfix.DateTimeSeconds:      "DateTime('Asia/Jakarta')",
	ipfix.DateTimeMilliseconds: "DateTime64(3, 'Asia/Jakarta')",
	ipfix.DateTimeMicroseconds: "DateTime64(6, 'Asia/Jakarta')",
	ipfix.DateTimeNanoseconds:  "DateTime64(9, 'Asia/Jakarta')",
	ipfix.Ipv4Address:          "IPv4",
	ipfix.Ipv6Address:          "IPv6",
}

var ipfixAvroMap = map[ipfix.FieldType][]string{
	ipfix.Uint8:                {"null", "int"},
	ipfix.Uint16:               {"null", "int"},
	ipfix.Uint32:               {"null", "long"},
	ipfix.Uint64:               {"null", "string"},
	ipfix.Int8:                 {"null", "int"},
	ipfix.Int16:                {"null", "int"},
	ipfix.Int32:                {"null", "int"},
	ipfix.Int64:                {"null", "long"},
	ipfix.Float32:              {"null", "float"},
	ipfix.Float64:              {"null", "double"},
	ipfix.Boolean:              {"null", "boolean"},
	ipfix.MacAddress:           {"null", "string"},
	ipfix.OctetArray:           {"null", "string"},
	ipfix.String:               {"null", "string"},
	ipfix.DateTimeSeconds:      {"null", "long"},
	ipfix.DateTimeMilliseconds: {"null", "string"},
	ipfix.DateTimeMicroseconds: {"null", "string"},
	ipfix.DateTimeNanoseconds:  {"null", "string"},
	ipfix.Ipv4Address:          {"null", "string"},
	ipfix.Ipv6Address:          {"null", "string"},
}

func clickhouseColumns() {
	for i := 1; i < 492; i++ {
		if val, ok := ipfix.InfoModel[ipfix.ElementKey{0, uint16(i)}]; ok {
			cName := stringy.New(val.Name)
			chType, ok := ipfixChMap[val.Type]
			if !ok {
				chType = "String"
			}
			fmt.Printf("ds_iana_%s Nullable(%s),\n", cName.SnakeCase().ToLower(), chType)
		}
	}
}

type avroschema struct {
	Name    string   `json:"name"`
	Type    []string `json:"type"`
	Aliases []string `json:"aliases"`
	Default any      `json:"default"`
}

func avroSchema() {
	var schema []avroschema
	for i := 1; i < 492; i++ {
		if val, ok := ipfix.InfoModel[ipfix.ElementKey{0, uint16(i)}]; ok {
			t, ok := ipfixAvroMap[val.Type]
			if !ok {
				t = []string{"null", "string"}
			}
			schema = append(schema, avroschema{
				Name:    val.Name,
				Type:    t,
				Default: nil,
				Aliases: []string{strconv.FormatUint(uint64(val.FieldID), 10)},
			})
		}
	}
	o, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(o))
}

func avroToChColumns() {
	for i := 1; i < 492; i++ {
		if val, ok := ipfix.InfoModel[ipfix.ElementKey{0, uint16(i)}]; ok {
			cName := stringy.New(val.Name)
			fmt.Printf("RPATH_STRING(DataSets[1], '/%s') as ds_iana_%s,\n", val.Name, cName.SnakeCase().ToLower())
		}
	}
}

func main() {
	avroToChColumns()
}
