package course

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/yiplee/blockquiz/core"
	"github.com/yiplee/blockquiz/store"
)

type fileLoader struct {
	courses []*core.Course
	set     map[int64]bool
}

func LoadCourses(courseFolder string) core.CourseStore {
	s := &fileLoader{
		set: make(map[int64]bool),
	}

	err := filepath.Walk(courseFolder, func(path string, info os.FileInfo, err error) error {
		ext := filepath.Ext(path)
		if ext != ".yaml" {
			return nil
		}

		name := strings.TrimSuffix(info.Name(), ext)
		fields := strings.Fields(name)

		id, _ := strconv.ParseInt(fields[0], 10, 64)
		if id <= 0 {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		var course core.Course
		if err := yaml.NewDecoder(f).Decode(&course); err != nil {
			return err
		}

		if err := core.ValidateCourse(&course); err != nil {
			return fmt.Errorf("validate %s failed: %w", course.Title, err)
		}

		course.ID = id
		if err := s.insert(&course); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	s.sort()
	return s
}

func (s *fileLoader) sort() {
	sort.Slice(s.courses, func(i, j int) bool {
		return s.courses[i].ID < s.courses[j].ID
	})
}

func (s *fileLoader) insert(course *core.Course) error {
	if s.set[course.ID] {
		return fmt.Errorf("dumplicated course %d inserted", course.ID)
	}

	s.courses = append(s.courses, course)
	s.set[course.ID] = true
	return nil
}

func (s *fileLoader) Add(ctx context.Context, course *core.Course) error {
	panic("adding course runtime is forbidden")
}

func (s *fileLoader) ListAll(ctx context.Context) ([]*core.Course, error) {
	return s.courses[:], nil
}

func (s *fileLoader) ListLanguage(ctx context.Context, language string) ([]*core.Course, error) {
	courses := make([]*core.Course, 0, len(s.courses))
	for _, course := range s.courses {
		if course.Language == language {
			courses = append(courses, course)
		}
	}

	return courses, nil
}

func (s *fileLoader) Find(ctx context.Context, id int64) (*core.Course, error) {
	courses := s.courses
	idx := sort.Search(len(courses), func(i int) bool {
		return courses[i].ID >= id
	})

	if idx >= len(courses) {
		return nil, store.ErrNotFound
	}

	course := courses[idx]
	if course.ID != id {
		return nil, store.ErrNotFound
	}

	return course, nil
}

func (s *fileLoader) FindNext(ctx context.Context, course *core.Course) (*core.Course, error) {
	courses, err := s.ListLanguage(ctx, course.Language)
	if err != nil {
		return nil, err
	}

	if len(courses) == 0 {
		return nil, store.ErrNotFound
	}

	idx := sort.Search(len(courses), func(i int) bool {
		c := courses[i]
		return c.ID > course.ID
	})

	if idx >= len(courses) {
		// get from beginning
		idx = 0
	}

	return courses[idx], nil
}
