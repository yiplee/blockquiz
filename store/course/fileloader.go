package course

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	jsoniter "github.com/json-iterator/go"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

type fileLoader struct {
	courses []json.RawMessage
	indexes map[string]int
}

func LoadCourses(courseFolder string) core.CourseStore {
	s := &fileLoader{
		courses: make([]json.RawMessage, 0),
		indexes: make(map[string]int),
	}

	err := filepath.Walk(courseFolder, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext != ".yaml" {
			return nil
		}

		name := strings.TrimSuffix(info.Name(), ext)
		fields := strings.Split(name, ".")

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var course core.Course
		if err := yaml.NewDecoder(f).Decode(&course); err != nil {
			return err
		}

		if course.Title == "" {
			course.Title = fields[0]
		}

		if course.Language == "" {
			course.Language = fields[1]
		}

		if err := core.ValidateCourse(&course); err != nil {
			return fmt.Errorf("validate %s failed: %w", course.Title, err)
		}

		if err := s.insert(&course); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	return s
}

func courseKey(title, language string) string {
	return title + language
}

func (s *fileLoader) insert(course *core.Course) error {
	key := courseKey(course.Title, course.Language)
	if _, ok := s.indexes[key]; ok {
		return fmt.Errorf("dumplicated course %s %s inserted", course.Title, course.Language)
	}

	data, _ := jsoniter.Marshal(course)
	s.courses = append(s.courses, data)
	s.indexes[key] = len(s.courses) - 1
	return nil
}

func (s *fileLoader) Add(ctx context.Context, course *core.Course) error {
	panic("adding course runtime is forbidden")
}

func (s *fileLoader) ListAll(ctx context.Context) ([]*core.Course, error) {
	buff := &bytes.Buffer{}
	if err := jsoniter.NewEncoder(buff).Encode(s.courses); err != nil {
		return nil, err
	}

	var courses []*core.Course
	if err := jsoniter.NewDecoder(buff).Decode(&courses); err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *fileLoader) ListLanguage(ctx context.Context, language string) ([]*core.Course, error) {
	courses, err := s.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	var idx int
	for _, course := range courses {
		if course.Language == language {
			courses[idx] = course
			idx++
		}
	}

	return courses[:idx], nil
}

func (s *fileLoader) Find(ctx context.Context, title, language string) (*core.Course, error) {
	key := courseKey(title, language)
	idx, ok := s.indexes[key]
	if !ok {
		return nil, store.ErrNotFound
	}

	data := s.courses[idx]
	var course core.Course
	err := jsoniter.Unmarshal(data, &course)
	return &course, err
}
