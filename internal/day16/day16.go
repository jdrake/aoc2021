package day16

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

var hexMap = map[string]string{
	"0": "0000",
	"1": "0001",
	"2": "0010",
	"3": "0011",
	"4": "0100",
	"5": "0101",
	"6": "0110",
	"7": "0111",
	"8": "1000",
	"9": "1001",
	"A": "1010",
	"B": "1011",
	"C": "1100",
	"D": "1101",
	"E": "1110",
	"F": "1111",
}

func ParseHex(line string) []int {
	s := ""
	for _, char := range strings.Split(line, "") {
		s += hexMap[char]
	}
	var bits []int
	for _, char := range strings.Split(s, "") {
		bit, _ := strconv.Atoi(char)
		bits = append(bits, bit)
	}
	return bits
}

func Bin2Int(s string) int {
	i, _ := strconv.ParseInt(s, 2, 64)
	return int(i)
}

func Bits2Int(bits []int) int {
	s := ""
	for _, b := range bits {
		s += strconv.Itoa(b)
	}
	return Bin2Int(s)
}

const Sum = 0
const Product = 1
const Minimum = 2
const Maximum = 3
const Literal = 4
const Greater = 5
const Less = 6
const Equal = 7

func CalcLiteral(bits []int) (int, int) {
	packetLength := 6
	var s string
	for {
		final := bits[packetLength] == 0
		packetLength += 1
		group := ""
		for j := 0; j < 4; j++ {
			group += strconv.Itoa(bits[packetLength+j])
		}
		fmt.Println("group", group)
		s += group
		packetLength += 4
		if final {
			break
		}
	}
	literal := Bin2Int(s)
	fmt.Println("literal", literal, "length", packetLength)
	return literal, packetLength
}

func GetSubpacketValues(bits []int) ([]int, int) {
	packetLength := 7
	var values []int
	if bits[6] == 0 {
		// Check next 15 for total length count
		packetLength += 15
		totalLength := Bits2Int(bits[7:22])
		start := 22
		for packetLength < totalLength+15 {
			value, subpacketLength := CalculatePacketValue(bits[start:])
			values = append(values, value)
			packetLength += subpacketLength
			start += subpacketLength
		}
	} else {
		// Check next 11 for subpacket count
		packetLength += 11
		subpacketCount := Bits2Int(bits[7:18])
		start := 18
		for i := 0; i < subpacketCount; i++ {
			value, subpacketLength := CalculatePacketValue(bits[start:])
			values = append(values, value)
			packetLength += subpacketLength
			start += subpacketLength
		}
	}
	fmt.Println("operator", bits[6], "values", values, "length", packetLength)
	return values, packetLength
}

func CalcSum(bits []int) (int, int) {
	sum := 0
	values, packetLength := GetSubpacketValues(bits)
	for _, value := range values {
		sum += value
	}
	return sum, packetLength
}

func CalcProduct(bits []int) (int, int) {
	product := 1
	values, packetLength := GetSubpacketValues(bits)
	for _, value := range values {
		product *= value
	}
	return product, packetLength
}

func CalcMinimum(bits []int) (int, int) {
	min := math.MaxInt
	values, packetLength := GetSubpacketValues(bits)
	for _, value := range values {
		if value < min {
			min = value
		}
	}
	return min, packetLength
}

func CalcMaximum(bits []int) (int, int) {
	max := 0
	values, packetLength := GetSubpacketValues(bits)
	for _, value := range values {
		if value > max {
			max = value
		}
	}
	return max, packetLength
}

func CalcGreater(bits []int) (int, int) {
	values, packetLength := GetSubpacketValues(bits)
	var value int
	if values[0] > values[1] {
		value = 1
	} else {
		value = 0
	}
	return value, packetLength
}

func CalcLess(bits []int) (int, int) {
	values, packetLength := GetSubpacketValues(bits)
	var value int
	if values[0] < values[1] {
		value = 1
	} else {
		value = 0
	}
	return value, packetLength
}

func CalcEqual(bits []int) (int, int) {
	values, packetLength := GetSubpacketValues(bits)
	var value int
	if values[0] == values[1] {
		value = 1
	} else {
		value = 0
	}
	return value, packetLength
}

func CalculatePacketValue(bits []int) (int, int) {
	packetType := Bits2Int(bits[3:6])
	switch packetType {
	case Sum:
		return CalcSum(bits)
	case Product:
		return CalcProduct(bits)
	case Minimum:
		return CalcMinimum(bits)
	case Maximum:
		return CalcMaximum(bits)
	case Literal:
		return CalcLiteral(bits)
	case Greater:
		return CalcGreater(bits)
	case Less:
		return CalcLess(bits)
	case Equal:
		return CalcEqual(bits)
	default:
		log.Fatal("unknown operator type:", packetType)
	}
	return -1, -1
}

func Day16() {
	// bits := ParseHex("C200B40A82")
	// bits := ParseHex("04005AC33890")
	// bits := ParseHex("880086C3E88112")
	// bits := ParseHex("CE00C43D881120")
	// bits := ParseHex("D8005AC2A8F0")
	// bits := ParseHex("F600BC2D8F")
	// bits := ParseHex("9C005AC2F8F0")
	// bits := ParseHex("9C0141080250320F1802104A08")
	bits := ParseHex("20546718027401204FE775D747A5AD3C3CCEEB24CC01CA4DFF2593378D645708A56D5BD704CC0110C469BEF2A4929689D1006AF600AC942B0BA0C942B0BA24F9DA8023377E5AC7535084BC6A4020D4C73DB78F005A52BBEEA441255B42995A300AA59C27086618A686E71240005A8C73D4CF0AC40169C739584BE2E40157D0025533770940695FE982486C802DD9DC56F9F07580291C64AAAC402435802E00087C1E8250440010A8C705A3ACA112001AF251B2C9009A92D8EBA6006A0200F4228F50E80010D8A7052280003AD31D658A9231AA34E50FC8010694089F41000C6A73F4EDFB6C9CC3E97AF5C61A10095FE00B80021B13E3D41600042E13C6E8912D4176002BE6B060001F74AE72C7314CEAD3AB14D184DE62EB03880208893C008042C91D8F9801726CEE00BCBDDEE3F18045348F34293E09329B24568014DCADB2DD33AEF66273DA45300567ED827A00B8657B2E42FD3795ECB90BF4C1C0289D0695A6B07F30B93ACB35FBFA6C2A007A01898005CD2801A60058013968048EB010D6803DE000E1C6006B00B9CC028D8008DC401DD9006146005980168009E1801B37E02200C9B0012A998BACB2EC8E3D0FC8262C1009D00008644F8510F0401B825182380803506A12421200CB677011E00AC8C6DA2E918DB454401976802F29AA324A6A8C12B3FD978004EB30076194278BE600C44289B05C8010B8FF1A6239802F3F0FFF7511D0056364B4B18B034BDFB7173004740111007230C5A8B6000874498E30A27BF92B3007A786A51027D7540209A04821279D41AA6B54C15CBB4CC3648E8325B490401CD4DAFE004D932792708F3D4F769E28500BE5AF4949766DC24BB5A2C4DC3FC3B9486A7A0D2008EA7B659A00B4B8ACA8D90056FA00ACBCAA272F2A8A4FB51802929D46A00D58401F8631863700021513219C11200996C01099FBBCE6285106")
	fmt.Println(bits)
	fmt.Println(CalculatePacketValue(bits))
}
