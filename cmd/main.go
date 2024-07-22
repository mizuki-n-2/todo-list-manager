package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
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

const dataFile = "data/todoList.json"

func main() {
	todoList := loadToDoList()
	scanner := bufio.NewScanner(os.Stdin)

	color.New(color.Bold).Println("Welcome to the TODO List Manager!")
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
				printWarnMsg("Usage: add <task name>")
				continue
			}
			taskName := strings.Join(args, " ")
			todoList.addTask(taskName)
		case "remove":
			if len(args) != 1 {
				printWarnMsg("Usage: remove <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				printErrMsg(fmt.Sprintf("Invalid task ID: %v", args[0]))
				continue
			}
			todoList.removeTask(taskID)
		case "list":
			if len(args) != 0 {
				printWarnMsg("Usage: list")
				continue
			}
			todoList.listTasks()
		case "complete":
			if len(args) != 1 {
				printWarnMsg("Usage: complete <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				printErrMsg(fmt.Sprintf("Invalid task ID: %v", args[0]))
				continue
			}
			todoList.completeTask(taskID)
		case "show":
			if len(args) != 1 {
				printWarnMsg("Usage: show <task ID>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				printErrMsg(fmt.Sprintf("Invalid task ID: %v", args[0]))
				continue
			}
			todoList.showTask(taskID)
		case "edit":
			if len(args) < 2 {
				printWarnMsg("Usage: edit <task ID> <new task name>")
				continue
			}
			taskID, err := strconv.Atoi(args[0])
			if err != nil {
				printErrMsg(fmt.Sprintf("Invalid task ID: %v", args[0]))
				continue
			}
			newTaskName := strings.Join(args[1:], " ")
			todoList.editTask(taskID, newTaskName)
		case "remove-all":
			if len(args) != 0 {
				printWarnMsg("Usage: remove-all")
				continue
			}
			todoList.removeAllTasks()
		case "help":
			fmt.Println("Commands:")
			fmt.Println("add <task name> - Add a new task")
			fmt.Println("remove <task ID> - Remove a task")
			fmt.Println("list - List all tasks")
			fmt.Println("complete <task ID> - Mark a task as completed")
			fmt.Println("show <task ID> - Show task details")
			fmt.Println("edit <task ID> <new task name> - Edit a task")
			fmt.Println("remove-all - Remove all tasks")
		case "exit":
			saveToDoList(todoList)
			fmt.Println("Save done...")
			fmt.Println("Goodbye!")
			return
		default:
			printErrMsg(fmt.Sprintf("Unknown command: %s", command))
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
	printInfoMsg(fmt.Sprintf("Task added: %+v", task))
}

// removeTask タスクの削除
func (t *ToDoList) removeTask(id int) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks = append(t.Tasks[:i], t.Tasks[i+1:]...)
			printInfoMsg(fmt.Sprintf("Task removed: %+v", task))
			return
		}
	}
	printErrMsg(fmt.Sprintf("Task not found with ID: %d", id))
}

// listTasks タスクの一覧表示
func (t *ToDoList) listTasks() {
	if len(t.Tasks) == 0 {
		printInfoMsg("No tasks")
		return
	}

	for _, task := range t.Tasks {
		status := "[ ]"
		if task.IsCompleted {
			status = "[x]"
		}
		printInfoMsg(fmt.Sprintf("%d. %s %s\n", task.ID, status, task.Name))
	}
}

// completeTask タスクの完了
func (t *ToDoList) completeTask(id int) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks[i].IsCompleted = true
			printInfoMsg(fmt.Sprintf("Task completed: %s", task.Name))
			return
		}
	}
	printErrMsg(fmt.Sprintf("Task not found with ID: %d", id))
}

// showTask タスクの詳細表示
func (t *ToDoList) showTask(id int) {
	for _, task := range t.Tasks {
		if task.ID == id {
			status := "Incomplete"
			if task.IsCompleted {
				status = "Complete"
			}
			printInfoMsg(fmt.Sprintf("Task ID: %d\nName: %s\nStatus: %s\n", task.ID, task.Name, status))
			return
		}
	}
	printErrMsg(fmt.Sprintf("Task not found with ID: %d", id))
}

// editTask タスクの編集
func (t *ToDoList) editTask(id int, newName string) {
	for i, task := range t.Tasks {
		if task.ID == id {
			t.Tasks[i].Name = newName
			printInfoMsg(fmt.Sprintf("Task edited: %s", newName))
			return
		}
	}
	printErrMsg(fmt.Sprintf("Task not found with ID: %d", id))
}

// removeAllTasks 全タスクの削除
func (t *ToDoList) removeAllTasks() {
	t.Tasks = []Task{}
	printInfoMsg("All tasks removed")
}

func loadToDoList() ToDoList {
	wd, _ := os.Getwd()
	file, err := os.Open(wd + "/" + dataFile)
	if err != nil {
		printErrMsg(fmt.Sprintf("Failed to open file: %v", err))
		return ToDoList{}
	}
	defer file.Close()

	var todoList ToDoList
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&todoList); err != nil {
		printErrMsg(fmt.Sprintf("Failed to decode JSON: %v", err))
		return ToDoList{}
	}
	return todoList
}

func saveToDoList(todoList ToDoList) {
	wd, _ := os.Getwd()
	file, err := os.Create(wd + "/" + dataFile)
	if err != nil {
		printErrMsg(fmt.Sprintf("Failed to create file: %v", err))
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(&todoList); err != nil {
		printErrMsg(fmt.Sprintf("Failed to encode JSON: %v", err))
	}
}

func printErrMsg(msg string) {
	color.New(color.FgRed).Println(msg)
}

func printWarnMsg(msg string) {
	color.New(color.FgYellow).Println(msg)
}

func printInfoMsg(msg string) {
	color.New(color.FgBlue).Println(msg)
}
