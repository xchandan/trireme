package packet

import "fmt"

type TCPOptions byte

const (
	EndOfOptionsList                  TCPOptions = 0
	NOP                                          = 1
	MSS                                          = 2
	WindowScale                                  = 3
	SackPermitted                                = 4
	Sack                                         = 5
	Echo                                         = 6
	EchoReply                                    = 7
	TimeStamps                                   = 8
	PartialOrderConnectionPermitted              = 9
	PartialOrderServiceProfile                   = 10
	CC                                           = 11
	CCNEW                                        = 12
	CCECHO                                       = 13
	TCPAlternateChecksumRequest                  = 14
	TCPAlternateChecksumData                     = 15
	Skeeter                                      = 16
	Bubba                                        = 17
	TrailerChecksumOption                        = 18
	MD5SignatureOption                           = 19
	SCPSCapabilities                             = 20
	SelectiveNegativeAcknowledgements            = 21
	RecordBoundaries                             = 22
	CorruptionExperienced                        = 23
	SNAP                                         = 24
	Unassigned                                   = 25
	TCPCompressionFilter                         = 26
	QuickStartResponse                           = 27
	UserTimeoutOption                            = 28
	TCPAuthenticationOption                      = 29
	MultipathTCP                                 = 30
	Reserved                                     = 31
	TCPFastopenCookie                            = 34
	ReservedRangeBegin                           = 35
	ReservedRangeEnd                             = 254
	AporetoAuthentication                        = 255
)

type tcpOptionsFormat struct {
	kind   TCPOptions
	length int
	data   []byte
}

// 0 - indicates no payload with the option
// -1 - variable length payload -- read from packet while parsing
// >0 - standard header data
var optionsMap = map[TCPOptions]int{
	EndOfOptionsList:                0,
	NOP:                             0,
	MSS:                             4,
	WindowScale:                     3,
	SackPermitted:                   2,
	Sack:                            -1,
	Echo:                            6,
	EchoReply:                       6,
	TimeStamps:                      10,
	PartialOrderConnectionPermitted: 2,
	PartialOrderServiceProfile:      3,
	CC:     6,
	CCNEW:  6,
	CCECHO: 6,
	TCPAlternateChecksumRequest: 3,
	TCPAlternateChecksumData:    -1,
	Skeeter:                     0,
	Bubba:                       0,
	TrailerChecksumOption:             3,
	MD5SignatureOption:                18,
	SCPSCapabilities:                  4,
	SelectiveNegativeAcknowledgements: -1,
	RecordBoundaries:                  2,
	CorruptionExperienced:             2,
	SNAP:                    -1,
	Unassigned:              0,
	TCPCompressionFilter:    0,
	QuickStartResponse:      8,
	UserTimeoutOption:       4,
	TCPAuthenticationOption: 0,
	MultipathTCP:            -1,
	Reserved:                0,
	TCPFastopenCookie:       -1,
	AporetoAuthentication:   4,
}

func (p *Packet) parseTCPOption(bytes []byte) {
	var index byte
	options := bytes[TCPOptionPos:p.TCPDataStartBytes()]
	for index < byte(len(options)) {
		if optionsMap[TCPOptions(options[index])] == 0 {

			p.L4TCPPacket.optionsMap[TCPOptions(options[index])] = tcpOptionsFormat{
				kind:   TCPOptions(options[index]),
				length: 0,
				data:   []byte{},
			}
			index = index + 1
		} else if optionsMap[TCPOptions(options[index])] == -1 {
			p.L4TCPPacket.optionsMap[TCPOptions(options[index])] = tcpOptionsFormat{
				kind:   TCPOptions(options[index]),
				length: int(options[index+1]),
				data:   options[index:(int(options[index+1]) - 2 + int(index))],
			}
			index = index + options[index+1]
		} else {

			p.L4TCPPacket.optionsMap[TCPOptions(options[index])] = tcpOptionsFormat{
				kind:   TCPOptions(options[index]),
				length: optionsMap[TCPOptions(options[index])],
				data:   options[index:(int(index) - 2 + optionsMap[TCPOptions(options[index])])],
			}
			index = index + byte(optionsMap[TCPOptions(options[index])])
		}

	}

}

//TCPOptionLength :: accessor function for option payload length
func (p *Packet) TCPOptionLength(option TCPOptions) int {
	return p.L4TCPPacket.optionsMap[option].length
}

//TCPOptionData :: accessor function to the slice of data
func (p *Packet) TCPOptionData(option TCPOptions) ([]byte, bool) {
	optionval, ok := p.L4TCPPacket.optionsMap[option]
	if ok {
		return optionval.data, true
	}

	return nil, false

}

//SetOptionData :: Rewrite data for an option that is already present
func (p *Packet) SetTCPOptionData(option TCPOptions, data []byte) {
	copy(p.L4TCPPacket.optionsMap[option].data, data)
}

//WalkTCPOption :: debug function
func (p *Packet) WalkTCPOptions() {
	fmt.Println("************************")
	for key, val := range p.L4TCPPacket.optionsMap {
		fmt.Println("$$$$$$$$$")
		fmt.Println(key)
		fmt.Println(val.length)
		fmt.Println(val.data)
		fmt.Println("$$$$$$$$$")
	}
	fmt.Println("************************")
}

func (p *Packet) TCPDataOffset() uint8 {
	return p.L4TCPPacket.tcpDataOffset
}
