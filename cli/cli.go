package cli

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func removeEmptyStrings(strArr []string) []string {
	var filteredArr []string

	for _, str := range strArr {
		if str != "" {
			filteredArr = append(filteredArr, str)
		}
	}
	return filteredArr
}

func ReadTodos() ([]byte, []string, error) {
	todoBytes, err := os.ReadFile("todos.txt")
	todos := removeEmptyStrings(strings.Split(string(todoBytes), "\n"))
	if err != nil {
		log.Fatalf("Failed to read todos: %v", err)
	}
	return todoBytes, todos, err
}

func listTodo() {
	_, todos, _ := ReadTodos()

	if len(todos) == 0 {
		fmt.Println("No todos has been added")
		return
	}

	for i, todo := range todos {
		if todo != "" {
			fmt.Printf("%d: %s\n", i+1, todo)
		}
	}

}

func createTodo(todo string) {
	todoFile, err := os.OpenFile("todos.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		log.Fatalf("failed to open todos file")
	}

	defer todoFile.Close()

	_, err = todoFile.WriteString(todo + "\n")

	if err != nil {
		log.Fatalf("Failed wiriting to file: %s", err)
	}

	fmt.Println("successfully created todo task")
	listTodo()

}

func updateTodo(todoIndex int, todoStr string) {
	todoIndex = todoIndex - 1
	_, todos, _ := ReadTodos()

	if todoIndex >= 0 && todoIndex < len(todos) {
		todos[todoIndex] = todoStr
	} else {
		log.Fatalf("Invalid todo number %d", todoIndex+1)
	}

	allTodos := strings.Join(removeEmptyStrings(todos), "\n")

	err := os.WriteFile("todos.txt", []byte(allTodos), 0644)

	if err != nil {
		log.Fatalf("Failed to update todos, %v", err)
	}

	fmt.Println("successfully updated todo")
	listTodo()
}

func deleteTodo(todoIndex int) {
	todoIndex = todoIndex - 1
	_, todos, _ := ReadTodos()
	todos = append(todos[:todoIndex], todos[todoIndex+1:]...)

	allTodos := strings.Join(todos, "\n")

	err := os.WriteFile("todos.txt", []byte(allTodos), 0644)
	if err != nil {
		log.Fatalf("Failed to delete todo, %v", err)
	}

	fmt.Println("successfully deleted todo")
	listTodo()

}

func help() {
	fmt.Println(`Usage: todo [command] [arguments]
	
	Commands:
	  list                List all todos
	  create -todo        Create a new todo
	  update -order -todo Update an existing todo
	  delete -order       Delete an existing todo
	
	Arguments:
	  -todo   Specify the todo task
	  -order  Specify the todo number
	
	Examples:
	  todo list
	  todo create -todo "Buy groceries"
	  todo update -order 1 -todo "Buy milk"
	  todo delete -order 1
	
	Run 'todo [command] -h' for more information on a command.`)
}

func Run() {

	commandMap := map[string]func(f *flag.FlagSet){
		"list": func(f *flag.FlagSet) {
			listTodo()
		},
		"create": func(f *flag.FlagSet) {
			if os.Args[2:3][0] != "-todo" {
				log.Fatalf("please use -todo flag when creating a todo task")
			}
			todoContent := f.String("todo", "", "Create Todo")
			f.Parse(os.Args[2:])
			createTodo(*todoContent)
		},
		"update": func(f *flag.FlagSet) {
			todoOrder := f.Int("order", -1, "Todo order")
			todoContent := f.String("todo", "", "Updated todo")
			f.Parse(os.Args[2:])
			updateTodo(*todoOrder, *todoContent)
		},
		"delete": func(f *flag.FlagSet) {
			todoOrder := f.Int("order", -1, "Todo order")
			f.Parse(os.Args[2:])
			deleteTodo(*todoOrder)
		},
		"help": func(f *flag.FlagSet) {
			help()
		},
	}

	if len(os.Args) < 2 {
		help()
		os.Exit(0)
	}

	command, ok := commandMap[os.Args[1]]

	if !ok {
		help()
		os.Exit(0)
	}

	flagSet := flag.NewFlagSet(os.Args[1], flag.ContinueOnError)
	command(flagSet)
}
