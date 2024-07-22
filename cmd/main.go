package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsCompleted bool   `json:"isCompleted"`
}

type ToDoList struct {
	Tasks []Task `json:"tasks"`
}

const dataFile = "todoList.json"

func main() {
	todoList := loadToDoList()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the TODO List Manager!")
	fmt.Println("Type 'exit' to quit.")

	for {
		fmt.Print("> ")

		scanner.Scan()
		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		command := parts[0]
		args := parts[1:]

		switch command {
		case "add":
			if len(args) < 1 {
				fmt.Println("Usage: add <task name>")
				continue
			}
			taskName := strings.Join(args, " ")
			todoList.addTask(taskName)
		case "remove":
			if len(args) != 1 {
				fmt.Println("Usage: remove <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID:", args[0])
				continue
			}
			todoList.removeTask(taskID)
		case "list":
			if len(args) != 0 {
				fmt.Println("Usage: list")
				continue
			}
			todoList.listTasks()
		case "complete":
			if len(args) != 1 {
				fmt.Println("Usage: complete <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID:", args[0])
				continue
			}
			todoList.completeTask(taskID)
		case "show":
			if len(args) != 1 {
				fmt.Println("Usage: show <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID:", args[0])
				continue
			}
			todoList.showTask(taskID)
		case "edit":
			if len(args) < 2 {
				fmt.Println("Usage: edit <task ID> <new task name>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID:", args[0])
				continue
			}
			newTaskName := strings.Join(args[1:], " ")
			todoList.editTask(taskID, newTaskName)
		case "remove-all":
			if len(args) != 0 {
				fmt.Println("Usage: remove-all")
				continue
			}
			todoList.removeAllTasks()
		case "exit":
			saveToDoList(todoList)
			fmt.Println("Save done...")
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command:", command)
		}
	}
}

// addTask タスクの追加
func (t *ToDoList) addTask(name string) {
	id := 1
	if len(t.Tasks) > 0 {
		id = t.Tasks[len(t.Tasks)-1].ID + 1
	}
	task := Task{
		ID:          id,
		Name:        name,
		IsCompleted: false,
	}
	t.Tasks = append(t.Tasks, task)
	fmt.Println("Task added:", task)
}

// removeTask タスクの削除
func (t *ToDoList) removeTask(id int) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			fmt.Println("Task removed:", task)
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}

// listTasks タスクの一覧表示
func (t *ToDoList) listTasks() {
	if len(t.Tasks) == 0 {
		fmt.Println("No tasks")
		return
	}

	for _, task := range t.Tasks {
		status := "[ ]"
		if task.IsCompleted {
			status = "[x]"
		}
		fmt.Printf("%d. %s %s\n", task.ID, status, task.Name)
	}
}

// completeTask タスクの完了
func (t *ToDoList) completeTask(id int) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks[i].IsCompleted = true
			fmt.Println("Task completed:", task)
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}

// showTask タスクの詳細表示
func (t *ToDoList) showTask(id int) {
	for _, task := range t.Tasks {
		if task.ID == id {
			status := "Incomplete"
			if task.IsCompleted {
				status = "Complete"
			}
			fmt.Printf("Task ID: %d\nName: %s\nStatus: %s\n", task.ID, task.Name, status)
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}

// editTask タスクの編集
func (t *ToDoList) editTask(id int, newName string) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks[i].Name = newName
			fmt.Println("Task edited:", newName)
			return
		}
	}
	fmt.Println("Task not found with ID:", id)
}

// removeAllTasks 全タスクの削除
func (t *ToDoList) removeAllTasks() {
	t.Tasks = []Task{}
	fmt.Println("All tasks removed")
}

func loadToDoList() ToDoList {
	wd, _ := os.Getwd()
	file, err := os.Open(wd + "/" + dataFile)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return ToDoList{}
	}
	defer file.Close()

	var todoList ToDoList
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&todoList); err != nil {
		fmt.Println("Failed to decode JSON:", err)
		return ToDoList{}
	}
	return todoList
}

func saveToDoList(todoList ToDoList) {
	wd, _ := os.Getwd()
	file, err := os.Create(wd + "/" + dataFile)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(&todoList); err != nil {
		fmt.Println("Failed to encode JSON:", err)
	}
}
