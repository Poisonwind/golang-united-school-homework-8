package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)


var (
	OPERATION string
	ITEM string
	FILENAME string
	ID string
)


type Arguments map[string]string

type User struct {

	Id string `json:"id"`
	Email string `json:"email"`
	Age uint `json:"age"`

}


func readUsers(fileName string)(users []User, err error){

	if _, err = os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		os.Create(fileName)
		return users, nil
	}

	bytesReaded, err := os.ReadFile(fileName)
	if err != nil {
		return users, fmt.Errorf("failed to read file, %w", err)
	}

	json.Unmarshal(bytesReaded, &users)

	return
}

func userExist(users []User, id string)(us User, ok bool){

	for _, user := range(users){
		if user.Id == id {
			return user, true
		}
	}

	return
}

func writeToFile(users []User, filename string)(err error){

	bytesToWrite, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users, %w", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file, %w", err)
	}
	defer file.Close()

	_, err = file.Write(bytesToWrite)
	if err != nil {
		return fmt.Errorf("failed to write file, %w", err)
	}

	return
}

func parseArgs() Arguments {

	flag.StringVar(&FILENAME, "fileName", FILENAME, "path to a file")
	flag.StringVar(&ITEM, "item", ITEM, "item")
	flag.StringVar(&OPERATION, "operation", OPERATION, "just text")
	flag.StringVar(&ID, "id", ID, "id")

	flag.Parse()	

	return Arguments{
		"id": ID,
		"item": ITEM,
		"operation": OPERATION,
		"fileName": FILENAME,
	}
}

func operationList(filename string, writer io.Writer)(err error ){

	bytesReaded, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("filed to read file, %w", err)
	}

	fmt.Fprint(writer, string(bytesReaded))
	return
}

func operationAdd(args Arguments, writer io.Writer) (err error){

	newUser := User{}
	err = json.Unmarshal([]byte(args["item"]), &newUser)
	if err != nil {
		return fmt.Errorf("failed to unmarshal %w", err)
	}

	users, err := readUsers(args["fileName"])
	if err != nil {
		return fmt.Errorf("failed to read users %w", err)
	}
	
	if _, ok := userExist(users, newUser.Id); ok {
		fmt.Fprintf(writer, "Item with id %s already exists", newUser.Id)
		return nil
	}

	users = append(users, newUser)

	err = writeToFile(users, args["fileName"])
	if err != nil{
		return err
	}

	return

}

func operationfindById(filename string, writer io.Writer, id string)(err error){

	users, err := readUsers(filename)
	if err != nil {
		return fmt.Errorf("failed to read users %w", err)
	}

	user, ok := userExist(users, id)
	if !ok {
		fmt.Fprint(writer, "")
		return
	}

	bytesToWrite, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("fail to marshal user, %w", err)
	}

	fmt.Fprint(writer, string(bytesToWrite))
	
	return
}

func operationRemove(filename string, writer io.Writer, id string)(err error){

	users, err := readUsers(filename)
	if err != nil {
		return fmt.Errorf("failed to read users %w", err)
	}

	_, ok := userExist(users, id);

	if !ok {

		fmt.Fprintf(writer, "Item with id %s not found", id)
		return
	}

	newUsers := []User{}

	for _, user := range(users){
		if user.Id != id {
			newUsers = append(newUsers, user)
		}
	}


	err = writeToFile(newUsers, filename)
	if err != nil{
		return err
	}

return

}



func Perform(args Arguments, writer io.Writer) (err error) {
	
	if args["operation"] == "" {
		err = fmt.Errorf("-operation flag has to be specified")
		return
	}

	if args["fileName"] == "" {
		err = fmt.Errorf("-fileName flag has to be specified")
		return
	}
	

	switch args["operation"] {

	case "add":

		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")			
		}

		err = operationAdd(args, writer)
		if err != nil {
			return err
		}
		
	case "list":

		err = operationList(args["fileName"], writer)
		if err != nil {
			return
		}

	case "findById":

		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")			
		}

		err = operationfindById(args["fileName"], writer, args["id"])
		if err != nil {
			return
		}

	case "remove":

		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")			
		}

		err = operationRemove(args["fileName"], writer, args["id"])
		if err != nil {
			return
		}


	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
	


	return
}

func main() {
	
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
	
}
