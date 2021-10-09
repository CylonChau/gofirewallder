package dbus

import (
	"reflect"
	"strings"
)

type Source struct {
	Address string `json:"address"`
	Mac     string `json:"mac"`
	Ipset   string `json:"ipset"`
	Invert  string `json:"invert"`
}

type Destination struct {
	Address string `json:"address"`
	Invert  string `json:"invert"`
}

type Service struct {
	Name string `json:"name"`
}

type Port struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}
type Protocol struct {
	Value string `json:"value"`
}

type IcmpBlock struct {
	Name string `json:"name"`
}
type IcmpType struct {
	Name string `json:"name"`
}

type ForwardPort struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	ToPort   string `json:"toport"`
	ToAddr   string `json:"toaddr"`
}

type Log struct {
	Prefix string `json:"prefix"`
	Level  string `json:"level"`
	Limit  Limit  `json:"limit"`
}
type Limit struct {
	Value string `json:"value"`
}
type Audit struct {
	Limit Limit `json:"limit"`
}
type Accept struct {
	Flag  bool
	Limit Limit `json:"limit"`
}
type Reject struct {
	Type  string `json:"type"`
	Limit Limit  `json:"limit"`
}
type Drop struct {
	Flag  bool
	Limit Limit `json:"limit"`
}

type Mark struct {
	Set   string `json:"set"`
	Limit Limit  `json:"limit"`
}

type SourcePort struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

type Rule struct {
	Family      string      `json:"family"`
	Source      Source      `json:"source"`
	Destination Destination `json:"destination"`
	Service     Service     `json:"service"`
	Port        Port        `json:"port"`
	Protocol    Protocol    `json:"protocol"`
	IcmpBlock   IcmpBlock   `json:"icmpblock"`
	IcmpType    IcmpType    `json:"icmptype"`
	ForwardPort ForwardPort `json:"forwardport"`
	Log         Log         `json:"log"`
	Audit       Audit       `json:"audit"`
	Accept      Accept      `json:"accept"`
	Reject      Reject      `json:"reject"`
	Drop        Drop        `json:"drop"`
	Mark        Mark        `json:"mark"`
}

type Interface struct {
	Name string `json:"name"`
}

/*
 * 对应firewalld zoneSettingd的顺序
   [
	   "", version
	   "", short
	   "", description
	   False,  Forward
	   DEFAULT_ZONE_TARGET,  target
	   [], service
	   [], port
	   [], icmp-blocks
	   False,  masquerade
	   [], forward-ports
	   [], interface
	   [], sources
	   [], rich
	   [], protocols
	   [], source-ports
	   False icmp-block-inversion
	]

*/

type Settings struct {
	Version            string        `json:"version"`
	Short              string        `json:"short"`
	Description        string        `json:"description"`
	Forward            bool          `json:"forward"`
	Targe              string        `json:"target"`
	Service            []string      `json:"service"`
	Port               []Port        `json:"porst"`
	IcmpBlock          []IcmpBlock   `json:"icmpblock"`
	Masquerade         bool          `json:"masquerade"`
	ForwardPort        []ForwardPort `json:"forwardport"`
	Interface          []Interface   `json:"interface"`
	Source             []Source      `json:"source"`
	Rule               []Rule        `json:"rule"`
	Protocol           []Protocol    `json:"protocol"`
	SourcePort         []SourcePort  `json:"sourceport"`
	IcmpBlockInversion bool          `json:"icmp-block-inversion""`
}

func (this *Source) IsEmpty() bool {
	return reflect.DeepEqual(this, &Source{})
}

func (this *Destination) IsEmpty() bool {
	return reflect.DeepEqual(this, &Destination{})
}

func (this *Service) IsEmpty() bool {
	return reflect.DeepEqual(this, &Service{})
}

func (this *Port) IsEmpty() bool {
	return reflect.DeepEqual(this, &Port{})
}

func (this *Protocol) IsEmpty() bool {
	return reflect.DeepEqual(this, &Protocol{})
}

func (this *IcmpBlock) IsEmpty() bool {
	return reflect.DeepEqual(this, &IcmpBlock{})
}

func (this *IcmpType) IsEmpty() bool {
	return reflect.DeepEqual(this, &IcmpType{})
}

func (this *Log) IsEmpty() bool {
	return reflect.DeepEqual(this, &Log{})
}

func (this *ForwardPort) IsEmpty() bool {
	return reflect.DeepEqual(this, &ForwardPort{})
}

func (this *Audit) IsEmpty() bool {
	return reflect.DeepEqual(this, &Audit{})
}

func (this *Accept) IsEmpty() bool {
	return reflect.DeepEqual(this, &Accept{})
}

func (this *Reject) IsEmpty() bool {
	return reflect.DeepEqual(this, &Reject{})
}

func (this *Drop) IsEmpty() bool {
	return reflect.DeepEqual(this, &Drop{})
}

func (this *Mark) IsEmpty() bool {
	return reflect.DeepEqual(this, &Mark{})
}

func (this *Limit) IsEmpty() bool {
	return reflect.DeepEqual(this, &Limit{})
}

func (this *Source) ToString() string {
	var str = " source "
	if this.Address != "" {
		str += "address=" + this.Address
	} else if this.Mac != "" {
		str += "mac=" + this.Mac
	} else {
		str += "ipset=" + this.Ipset
	}
	if this.Invert != "" {
		str += " "
		str += "invert=" + this.Invert
	}
	str += " "
	return str
}

func (this *Destination) ToString() string {
	var str = " destination "
	if this.Address != "" {
		str += "address=" + this.Address
	}
	if this.Invert != "" {
		str += " "
		str += "invert=" + this.Invert
	}
	str += " "
	return str
}

func (this *Service) ToString() string {
	var str = "service "
	if this.Name != "" {
		str += "name=" + this.Name
	}
	str += " "
	return str
}

func (this *Port) ToString() string {
	var str = "port "
	if this.Port != "" {
		str += "name=" + this.Port
	}
	if this.Protocol != "" {
		str += "protocol=" + this.Protocol
	}

	str += " "
	return str
}

func (this *Protocol) ToString() string {
	var str = "Protocol "
	if this.Value != "" {
		str += "value=" + this.Value
	}

	str += " "
	return str
}

func (this *IcmpBlock) ToString() string {
	var str = "icmp-block "
	if this.Name != "" {
		str += "name=" + this.Name
	}

	str += " "
	return str
}

func (this *IcmpType) ToString() string {
	var str = "icmp-type "
	if this.Name != "" {
		str += "name=" + this.Name
	}

	str += " "
	return str
}

func (this *ForwardPort) ToString() string {
	var str = "forward-port "

	if this.Port != "" {
		str += "port=" + this.Port
	}

	if this.Protocol != "" {
		str += " "
		str += "protocol=" + this.Protocol
	}

	if this.ToPort != "" {
		str += " "
		str += "to-port=" + this.ToPort
	}

	if this.ToAddr != "" {
		str += " "
		str += "to-addr=" + this.ToAddr
	}

	str += " "
	return str
}

func (this *Log) ToString() string {
	var str = "log"

	if this.Prefix != "" {
		str += " " + "prefix=" + this.Prefix
	}

	if this.Level != "" {
		str += " " + "level=" + this.Level
	}

	if !this.Limit.IsEmpty() {
		str += " " + "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Audit) ToString() string {
	var str = "audit"

	if !this.Limit.IsEmpty() {
		str += " " + "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Accept) ToString() string {
	var str string

	if this.Flag {
		str = "accept "
	}
	if !this.Limit.IsEmpty() {
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Reject) ToString() string {
	var str = "reject "

	if this.Type != "" {
		str += "type=" + this.Type
	}

	if !this.Limit.IsEmpty() {
		str += " "
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Drop) ToString() string {
	var str string

	if this.Flag {
		str = "drop "
	}
	if !this.Limit.IsEmpty() {
		str += "limit value=" + this.Limit.Value
	}
	str += " "
	return str
}

func (this *Mark) ToString() string {
	var str = "mark"

	if this.Set != "" {
		str += " "
		str += "set=" + this.Set
	}

	if !this.Limit.IsEmpty() {
		str += " "
		str += "limit value=" + this.Limit.Value
	}

	str += " "
	return str
}

func (this *Rule) ToString() (ruleString string) {
	ruleString = "rule "
	if this.Family != "" {
		ruleString += "family=" + this.Family
	}

	if !this.Source.IsEmpty() {
		ruleString += this.Source.ToString()
	}

	if !this.Destination.IsEmpty() {
		ruleString += this.Destination.ToString()
	}

	if !this.Service.IsEmpty() {
		ruleString += this.Service.ToString()
	}

	if !this.Port.IsEmpty() {
		ruleString += this.Port.ToString()
	}

	if !this.Protocol.IsEmpty() {
		ruleString += this.Protocol.ToString()
	}

	if !this.IcmpBlock.IsEmpty() {
		ruleString += this.IcmpBlock.ToString()
	}

	if !this.IcmpType.IsEmpty() {
		ruleString += this.IcmpType.ToString()
	}

	if !this.ForwardPort.IsEmpty() {
		ruleString += this.ForwardPort.ToString()
	}

	if !this.Log.IsEmpty() {
		ruleString += this.Log.ToString()
	}

	if !this.Audit.IsEmpty() {
		ruleString += this.Audit.ToString()
	}

	if !this.Accept.IsEmpty() {
		ruleString += this.Accept.ToString()
	}

	if !this.Reject.IsEmpty() {
		ruleString += this.Reject.ToString()
	}

	if !this.Drop.IsEmpty() {
		ruleString += this.Drop.ToString()
	}

	if !this.Mark.IsEmpty() {
		ruleString += this.Mark.ToString()
	}
	return
}

func stringToReject(slice []string) (reject Reject, ruleSlice []string) {
Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[1] {
		case "type":
			slice = removeSliceElement(slice, index)
			reject.Type = slice[index]
			slice = removeSliceElement(slice, index)
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			reject.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return reject, ruleSlice
}

func stringToMark(slice []string) (mark Mark, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "set":
			slice = removeSliceElement(slice, index)
			mark.Set = tmp_slice[1]
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			mark.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func stringToForwardPort(slice []string) (forwardPort ForwardPort, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "port":
			slice = removeSliceElement(slice, index)
			forwardPort.Port = tmp_slice[1]
			goto Label
		case "protocol":
			slice = removeSliceElement(slice, index)
			forwardPort.Protocol = tmp_slice[1]
			goto Label
		case "to-port":
			slice = removeSliceElement(slice, index)
			forwardPort.ToPort = tmp_slice[1]
			goto Label
		case "to-addr":
			slice = removeSliceElement(slice, index)
			forwardPort.ToAddr = tmp_slice[1]
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func stringToLog(slice []string) (log Log, ruleSlice []string) {

Label:
	for index, value := range slice {
		tmp_slice := strings.Split(value, "=")
		switch tmp_slice[0] {
		case "prefix":
			slice = removeSliceElement(slice, index)
			log.Prefix = tmp_slice[1]
			goto Label
		case "level":
			slice = removeSliceElement(slice, index)
			log.Level = tmp_slice[1]
			goto Label
		case "limit":
			slice = removeSliceElement(slice, index)
			tmp_slice := strings.Split(slice[index], "=")
			log.Limit = Limit{Value: tmp_slice[1]}
			slice = removeSliceElement(slice, index)
			goto Label
		}
	}
	ruleSlice = slice
	return
}

func StringToRule(str string) (rule *Rule) {

	strslice := strings.Split(str, " ")
	rule = &Rule{}
Label:
	for index, value := range strslice {
		switch value {
		case "rule":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Family = tmp_str[1]
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "source":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			source := Source{}
			switch tmp_str[0] {
			case "address":
				source.Address = tmp_str[1]
			case "mac":
				source.Mac = tmp_str[1]
			case "ipset":
				source.Ipset = tmp_str[1]
			case "invert":
				source.Invert = tmp_str[1]
			}
			rule.Source = source
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "destination":
			strslice = removeSliceElement(strslice, index)
			dst := Destination{}
			tmp_str := strings.Split(strslice[index], "=")
			switch tmp_str[0] {
			case "address":
				dst.Address = tmp_str[1]
			case "invert":
				dst.Invert = tmp_str[1]
			}
			rule.Destination = dst
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "service":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Service = Service{Name: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "port":
			strslice = removeSliceElement(strslice, index)
			port := strings.Split(strslice[index], "=")
			protocol := strings.Split(strslice[index+1], "=")
			rule.Port = Port{
				Port:     port[1],
				Protocol: protocol[1],
			}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "protocol":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.Protocol = Protocol{Value: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "icmp-block":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.IcmpBlock = IcmpBlock{Name: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "icmp-type":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index+1], "=")
			rule.IcmpType = IcmpType{Name: tmp_str[1]}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "forward-port":
			strslice = removeSliceElement(strslice, index)
			rule.ForwardPort, strslice = stringToForwardPort(strslice)
			goto Label
		case "log":
			strslice = removeSliceElement(strslice, index)
			rule.Log, strslice = stringToLog(strslice)
			goto Label
		case "audit":
			strslice = removeSliceElement(strslice, index)
			tmp_str := strings.Split(strslice[index], "=")
			rule.Audit = Audit{Limit: Limit{Value: tmp_str[1]}}
			strslice = removeSliceElement(strslice, index)
			goto Label
		case "accept":
			strslice = removeSliceElement(strslice, index)
			rule.Accept.Flag = true
			var tmp_str []string
			if len(strslice) > 0 {
				if strslice[index] == "limit" {
					strslice = removeSliceElement(strslice, index)
					tmp_str = strings.Split(strslice[index], "=")
					rule.Accept.Limit = Limit{Value: tmp_str[1]}
				}
			}
			goto Label
		case "drop":
			var tmp_str []string
			strslice = removeSliceElement(strslice, index)
			rule.Drop.Flag = true
			if len(strslice) > 0 {
				if strslice[index] == "limit" {
					strslice = removeSliceElement(strslice, index)
					tmp_str = strings.Split(strslice[index], "=")
					rule.Drop.Limit = Limit{Value: tmp_str[1]}
				}
			}
			goto Label
		case "reject":
			strslice = removeSliceElement(strslice, index)
			rule.Reject, strslice = stringToReject(strslice)
			goto Label
		case "mark":
			strslice = removeSliceElement(strslice, index)
			rule.Mark, strslice = stringToMark(strslice)
			goto Label
		}
	}
	return rule
}
