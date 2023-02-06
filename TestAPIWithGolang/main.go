package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Student struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Old uint `json:"old`
	School string `json:"school`
}

var students = []Student{
	{ID: "1", Name: "A", Old: 18, School: "A School"},
	{ID: "2", Name: "B", Old: 18, School: "B School"},
	{ID: "3", Name: "C", Old: 18, School: "C School"},
}

func getStudentsList(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, students)
}

func postStudentList(c *gin.Context) {
	var newStudent Student

	if err := c.Bind(&newStudent); err != nil {
		return
	}

	students = append(students, newStudent)
	c.IndentedJSON(http.StatusCreated, students)
}

func getStudentByID(c *gin.Context) {
	id := c.Param("id")

	for _, student := range students {
		if student.ID == id {
			c.IndentedJSON(http.StatusOK, student)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, map[string]string{"message": "Student not found"})
}

func main() {
	router := gin.Default()

	router.GET("/student", getStudentsList)
	router.POST("/student", postStudentList)

	router.GET("/student/:id", getStudentByID)

	router.Run("localhost:8080")

}