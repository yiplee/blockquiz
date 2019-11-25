package course

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

type fileLoader struct {
	courses []core.Course
	indexes map[string]int
}

func LoadCourses(courseFolder string) core.CourseStore {
	s := &fileLoader{
		courses: make([]core.Course, 0),
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

		if err := s.insert(course); err != nil {
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

func (s *fileLoader) insert(course core.Course) error {
	key := courseKey(course.Title, course.Language)
	if _, ok := s.indexes[key]; ok {
		return fmt.Errorf("dumplicated course %s %s inserted", course.Title, course.Language)
	}

	s.courses = append(s.courses, course)
	s.indexes[key] = len(s.courses) - 1
	return nil
}

func (s *fileLoader) Add(ctx context.Context, course *core.Course) error {
	panic("adding course runtime is forbidden")
}

func (s *fileLoader) ListAll(ctx context.Context) ([]*core.Course, error) {
	courses := make([]*core.Course, 0, len(s.courses))
	for _, course := range s.courses {
		course := course
		courses = append(courses, &course)
	}

	return courses, nil
}

func (s *fileLoader) ListLanguage(ctx context.Context, language string) ([]*core.Course, error) {
	courses := make([]*core.Course, 0, len(s.courses))
	for _, course := range s.courses {
		if course := course; course.Language == language {
			courses = append(courses, &course)
		}
	}

	return courses, nil
}

func (s *fileLoader) Find(ctx context.Context, title, language string) (*core.Course, error) {
	key := courseKey(title, language)
	idx, ok := s.indexes[key]
	if !ok {
		return nil, store.ErrNotFound
	}

	var course core.Course
	course = s.courses[idx]
	return &course, nil
}
