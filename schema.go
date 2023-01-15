package web

import (
	"fmt"

	"lib.dev/english"
)

type Schema struct {
	Fields []Field
}

func (s *Schema) AddField(name english.Name, t Type) error {
	for _, f := range s.Fields {
		if f.EnglishName.String() == name.String() {
			return fmt.Errorf("field %s already exists", name)
		}
	}
	s.Fields = append(s.Fields, Field{
		EnglishName: name,
		Type:        t,
	})
	return nil
}

func (s *Schema) RemoveField(name english.Name) error {
	for i, f := range s.Fields {
		if f.EnglishName.String() == name.String() {
			s.Fields = append(s.Fields[:i], s.Fields[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("field %s does not exist", name)
}

func (s *Schema) MoveField(fromIndex, toIndex int) error {
	if fromIndex < 0 || fromIndex >= len(s.Fields) {
		return fmt.Errorf("invalid fromIndex %d", fromIndex)
	}
	if toIndex < 0 || toIndex >= len(s.Fields) {
		return fmt.Errorf("invalid toIndex %d", toIndex)
	}
	field := s.Fields[fromIndex]
	s.Fields = append(s.Fields[:fromIndex], s.Fields[fromIndex+1:]...)
	s.Fields = append(s.Fields[:toIndex], append([]Field{field}, s.Fields[toIndex:]...)...)
	return nil
}

func (s *Schema) ChangeFieldName(oldName, newName english.Name) error {
	for _, f := range s.Fields {
		if f.EnglishName.String() == newName.String() {
			return fmt.Errorf("field %s already exists", newName)
		}
	}
	for i, f := range s.Fields {
		if f.EnglishName.String() == oldName.String() {
			s.Fields[i].EnglishName = newName
			return nil
		}
	}
	return fmt.Errorf("field %s does not exist", oldName)
}

func (s *Schema) ChangeFieldType(fieldName english.Name, newType Type) error {
	for i, f := range s.Fields {
		if f.EnglishName.String() == fieldName.String() {
			s.Fields[i].Type = newType
			return nil
		}
	}
	return fmt.Errorf("field %s does not exist", fieldName)
}
