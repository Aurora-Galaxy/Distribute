package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (stu Student) Average() float32 {
	var result float32
	for _, g := range stu.Grades {
		result += g.Score
	}
	return result / float32(len(stu.Grades))
}

type Students []Student

var (
	students      Students
	studentsMutex = new(sync.Mutex)
)

// 使用 for range 查找内容需要返回地址时一定注意，其取出的value是单独开辟的内存空间存放
// 1.22 版本后 可忽略
func (ss Students) GetById(id int) (*Student, error) {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("Student with Id %d not found", id)
}

type GradeType string

const (
	GradeQuiz GradeType = "Quiz"
	GradeTest GradeType = "Test"
	GradeExam GradeType = "Exam"
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
