package util

import (
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"time"
)

type IniParser struct {
	confReader *ini.File // config reader
}

type IniParserError struct {
	errorInfo string
}

func (e *IniParserError) Error() string { return e.errorInfo }

func (p *IniParser) Load(configFile string) error {
	conf, err := ini.Load(configFile)
	if err != nil {
		p.confReader = nil
		return err
	}
	p.confReader = conf
	return nil
}

func (p *IniParser) GetString(section string, key string) string {
	s, _ := p.getSection(section)
	if s == nil {
		return ""
	}
	return s.Key(key).String()
}

func (p *IniParser) GetInt32(section string, key string) int32 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueInt, _ := s.Key(key).Int()
	return int32(valueInt)
}

func (p *IniParser) GetUint32(section string, key string) uint32 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueInt, _ := s.Key(key).Uint()
	return uint32(valueInt)
}

func (p *IniParser) GetInt64(section string, key string) int64 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueInt, _ := s.Key(key).Int64()
	return valueInt
}

func (p *IniParser) GetUint64(section string, key string) uint64 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueInt, _ := s.Key(key).Uint64()
	return valueInt
}

func (p *IniParser) GetFloat32(section string, key string) float32 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueFloat, _ := s.Key(key).Float64()
	return float32(valueFloat)
}

func (p *IniParser) GetFloat64(section string, key string) float64 {
	s, _ := p.getSection(section)
	if s == nil {
		return 0
	}
	valueFloat, _ := s.Key(key).Float64()
	return valueFloat
}

func (p *IniParser) GetBool(section string, key string) bool {
	s, _ := p.getSection(section)
	if s == nil {
		return false
	}
	valueBool, _ := s.Key(key).Bool()
	return valueBool
}

func (p *IniParser) GetDuration(section string, key string) time.Duration {
	s, _ := p.getSection(section)
	if s == nil {
		return time.Duration(0)
	}
	valueDuration, _ := s.Key(key).Duration()
	return valueDuration
}

func (p *IniParser) getSection(section string) (*ini.Section, error) {
	if section == "" {
		section = ini.DefaultSection
	}
	if p.confReader == nil {
		fmt.Errorf("config init failed, section: %v", section)
		return nil, errors.New("config init failed, section: " + section)
	}
	s := p.confReader.Section(section)
	if s == nil {
		fmt.Errorf("config read failed, section: %v", section)
		return nil, errors.New("config read failed, section: " + section)
	}
	return s, nil
}
