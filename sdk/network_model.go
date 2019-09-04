// Copyright 2018 ProximaX Limited. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package sdk

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/proximax-storage/go-xpx-utils/str"
	"io"
	"sort"
	"strings"
)

type NetworkType uint8

const (
	Mijin           NetworkType = 96
	MijinTest       NetworkType = 144
	Public          NetworkType = 184
	PublicTest      NetworkType = 168
	Private         NetworkType = 200
	PrivateTest     NetworkType = 176
	NotSupportedNet NetworkType = 0
	AliasAddress    NetworkType = 145
)

func NetworkTypeFromString(networkType string) NetworkType {
	switch networkType {
	case "mijin":
		return Mijin
	case "mijinTest":
		return MijinTest
	case "public":
		return Public
	case "publicTest":
		return PublicTest
	case "private":
		return Private
	case "privateTest":
		return PrivateTest
	}

	return NotSupportedNet
}

func (nt NetworkType) String() string {
	return fmt.Sprintf("%d", nt)
}

var networkTypeError = errors.New("wrong raw NetworkType value")

func ExtractNetworkType(version int64) NetworkType {
	return NetworkType(uint32(version) >> 24)
}

type Entity struct {
	Name              string
	Type              EntityType
	SupportedVersions []EntityVersion
}

func (e Entity) String() string {
	return str.StructToString(
		"Entity",
		str.NewField("Name", str.StringPattern, e.Name),
		str.NewField("Type", str.StringPattern, e.Type),
		str.NewField("SupportedVersions", str.StringPattern, e.SupportedVersions),
	)
}

type Field struct {
	Key     string
	Value   string
	Comment string
	Index   int
}

func NewField() *Field {
	return &Field{}
}

func (c Field) String() string {
	s := c.Comment
	if len(strings.TrimSpace(s)) != 0 {
		s += "\n"
	}

	s += fmt.Sprintf("%s = %s", c.Key, c.Value)
	return s
}

type ConfigBag struct {
	Name    string
	Comment string
	Index   int
	Fields  map[string]*Field
}

func NewConfigBag() *ConfigBag {
	return &ConfigBag{
		Fields: make(map[string]*Field),
	}
}

func (c ConfigBag) String() string {
	s := c.Comment
	if len(strings.TrimSpace(s)) != 0 {
		s += "\n"
	}

	s += fmt.Sprintf("[%s]\n", c.Name)

	fields := make([]*Field, 0, len(c.Fields))
	for _, f := range c.Fields {
		fields = append(fields, f)
	}

	sort.Slice(fields, func(i, j int) bool {
		return fields[i].Index < fields[j].Index
	})

	for _, field := range fields {
		s += field.String()
		s += "\n"
	}

	return s
}

type NetworkConfig struct {
	Sections map[string]*ConfigBag
}

func NewNetworkConfig() *NetworkConfig {
	c := NetworkConfig{
		Sections: make(map[string]*ConfigBag),
	}
	return &c
}

func (c *NetworkConfig) UnmarshalBinary(data []byte) error {
	const HASH = '#'
	const SEMICOLON = ';'
	const L_BRACKET = '['
	const R_BRACKET = ']'

	r := bufio.NewReader(bytes.NewReader(data))
	l := 0
	var bag *ConfigBag = nil
	comment := ""

	for true {
		line, isPrefix, err := r.ReadLine()
		l += 1

		if isPrefix {
			return fmt.Errorf("Line %d is to long", l)
		}

		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("Got error during read the line %d", l)
			}
			break
		}

		lineS := strings.TrimSpace(string(line))

		if len(lineS) == 0 {
			comment += "\n"
			continue
		}

		switch lineS[0] {
		case SEMICOLON, HASH:
			comment += lineS
		case L_BRACKET:
			bag = NewConfigBag()
			bag.Comment = comment
			bag.Index = len(c.Sections)
			comment = ""

			left := 1
			right := strings.Index(lineS, string(R_BRACKET))

			if right == -1 {
				return fmt.Errorf("Wrong header of section at line %d", l)
			}

			bag.Name = strings.TrimSpace(lineS[left:right])

			if _, ok := c.Sections[bag.Name]; ok {
				return fmt.Errorf("Duplicate section at line %d with name %s", l, bag.Name)
			}

			c.Sections[bag.Name] = bag
		default:
			separatorIndex := strings.Index(lineS, "=")

			switch separatorIndex {
			case -1:
				return fmt.Errorf("'=' character not found at line %d", l)
			case 0:
				return fmt.Errorf("Key is empty at line %d", l)
			default:
				if bag == nil {
					return fmt.Errorf("The section without header at line %d", l)
				}

				field := NewField()

				field.Key = strings.TrimSpace(lineS[0:separatorIndex])
				field.Value = strings.TrimSpace(lineS[separatorIndex+1:])
				field.Comment = comment
				field.Index = len(bag.Fields)
				comment = ""

				bag.Fields[field.Key] = field
			}
		}
	}

	return nil
}

func (c *NetworkConfig) MarshalBinary() (data []byte, err error) {
	s := ""

	sections := make([]*ConfigBag, 0, len(c.Sections))
	for _, f := range c.Sections {
		sections = append(sections, f)
	}

	sort.Slice(sections, func(i, j int) bool {
		return sections[i].Index < sections[j].Index
	})

	for _, section := range sections {
		s += section.String()
	}

	return []byte(s), nil
}

func (c *NetworkConfig) String() string {
	s, _ := c.MarshalBinary()

	return string(s)
}

type SupportedEntities struct {
	Entities map[EntityType]*Entity
}

func NewSupportedEntities() *SupportedEntities {
	ref := &SupportedEntities{
		Entities: make(map[EntityType]*Entity),
	}

	return ref
}

func (s *SupportedEntities) UnmarshalBinary(data []byte) error {
	dto := supportedEntitiesDTO{}
	err := json.Unmarshal(data, &dto)
	if err != nil {
		return err
	}

	err = dto.toStruct(s)
	if err != nil {
		return err
	}

	return nil
}

func (s *SupportedEntities) MarshalBinary() (data []byte, err error) {
	dto := supportedEntitiesDTO{
		Entities: make([]*entityDTO, len(s.Entities)),
	}

	i := 0
	for _, entity := range s.Entities {
		dto.Entities[i] = &entityDTO{
			Name:              entity.Name,
			Type:              fmt.Sprintf("%d", entity.Type),
			SupportedVersions: entity.SupportedVersions,
		}

		i += 1
	}

	sort.Slice(dto.Entities, func(i, j int) bool {
		return dto.Entities[i].Name < dto.Entities[j].Name
	})

	return json.MarshalIndent(&dto, "", "    ")
}

func (s *SupportedEntities) String() string {
	b, err := s.MarshalBinary()
	// We can't get error, because we created dto by self
	if err != nil {
		panic(err)
	}

	return string(b)
}

type BlockchainConfig struct {
	StartedHeight           Height
	NetworkConfig           *NetworkConfig
	SupportedEntityVersions *SupportedEntities
}

func (nc BlockchainConfig) String() string {
	return str.StructToString(
		"NetworkConfig",
		str.NewField("StartedHeight", str.StringPattern, nc.StartedHeight),
		str.NewField("NetworkConfig", str.StringPattern, nc.NetworkConfig),
		str.NewField("SupportedEntityVersions", str.StringPattern, nc.SupportedEntityVersions),
	)
}

type NetworkVersion struct {
	StartedHeight     Height
	BlockChainVersion BlockChainVersion
}

func (nv NetworkVersion) String() string {
	return str.StructToString(
		"NetworkVersion",
		str.NewField("StartedHeight", str.StringPattern, nv.StartedHeight),
		str.NewField("BlockChainVersion", str.StringPattern, nv.BlockChainVersion),
	)
}
