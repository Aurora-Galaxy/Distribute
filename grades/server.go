package grades

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func RegisterHandlers() {
	//handler := new(studentsHandler)
	http.Handle("/students", studentsHandler{})
	http.Handle("/students/", studentsHandler{})
}

type studentsHandler struct{}

//	/students
//
// /students/{id}
// /students/{id}/grades
// 正常处理可以使用正则表达式或对应的包接收url的参数，此处使用简单逻辑，作为演示使用
func (sh studentsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	//fmt.Println("splitpath = ", len(splitPath))
	switch len(splitPath) {
	case 2:
		//fmt.Println("===========")
		sh.getAll(w, r)
	case 3:
		id, err := strconv.Atoi(splitPath[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		sh.getOne(w, r, id)
	case 4:
		id, err := strconv.Atoi(splitPath[2])
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		fmt.Println("==========")
		sh.addGrade(w, r, id)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (sh studentsHandler) getAll(w http.ResponseWriter, r *http.Request) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()
	data, err := sh.toJson(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Failed to serialize students : %s", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (sh studentsHandler) getOne(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	// 根据 id 获取对应的内容
	student, err := students.GetById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Failed to serialize students : %s", err)
		return
	}

	data, err := sh.toJson(student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func (sh studentsHandler) addGrade(w http.ResponseWriter, r *http.Request, id int) {
	studentsMutex.Lock()
	defer studentsMutex.Unlock()

	// 根据 id 获取对应的内容
	student, err := students.GetById(id)
	fmt.Println("student = ", student)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Printf("Failed to serialize students : %s", err)
		return
	}

	var g Grade
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&g)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	//fmt.Println("add g = ", g)
	student.Grades = append(student.Grades, g)
	//fmt.Println("new  student = ", student)
	w.WriteHeader(http.StatusCreated)

	data, err := sh.toJson(g)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(data)
}

/**
 * toJson
 * @Description: 将传入的对象转换为字节切片
 * @receiver sh
 * @param obj
 * @return []byte
 * @return error
 */
func (sh studentsHandler) toJson(obj interface{}) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	err := enc.Encode(obj)
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize students : %s", err)
	}
	return b.Bytes(), nil
}
