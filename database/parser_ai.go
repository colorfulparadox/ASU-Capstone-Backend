package database

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
)

func New_AI(name string) error {
	var ai AI
	ai.Name = name
	return create_ai(ai)

}

func Add_Menu_Test(ai_name, menu_data string) string {
	menu_name, err := create_menu_file(menu_data, ai_name+"_menu", Menu_Path)
	if err != nil {
		log.Println(err)
		return "error"
	}

	return menu_name
}

func Add_Menu(ai_name, menu_data string) error {
	menu_name, err := create_menu_file(menu_data, ai_name+"_menu", Menu_Path)
	if err != nil {
		log.Println(err)
		return Invalid_Data()
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	file, err := os.Open(Menu_Path + menu_name)
	if err != nil {
		log.Println(err)
		return Invalid_Data()
	}
	defer file.Close()

	//Need to get old menu and delete it========================================================================
	//Or make a menu table so that the menu can be better catigorized also giving the ability for different menus

	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		log.Println(err)
		return Invalid_Data()
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		log.Println("File not found: ", err)
		return Invalid_Data()
	}
	log.Println(bs)

	var file_upload openai.FileBytesRequest

	file_upload.Name = menu_name
	file_upload.Bytes = bs
	file_upload.Purpose = "assistants"

	file_return, err := client.CreateFileBytes(context.Background(), file_upload)
	if err != nil {
		log.Println("File not stored in openai: ", err)
		return Invalid_Data()
	}

	var vector_upload openai.VectorStoreRequest

	vector_upload.Name = "menu"
	vector_upload.FileIDs = append(vector_upload.FileIDs, file_return.ID)

	vector_return, err := client.CreateVectorStore(context.Background(), vector_upload)
	if err != nil {
		log.Println("File not transfered to vector: ", err)
		return Invalid_Data()
	}

	ai, err := retrieve_ai(ai_name)
	if err != nil {
		log.Println("Error retrieveing the ai:", err)
		return Invalid_Data()
	}
	ai.FileID = file_return.ID
	ai.VectorID = vector_return.ID

	update_ai(ai)

	return nil
}
