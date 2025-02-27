package database

// type AI struct {
// 	AI_ID        int       `json:"id"`
// 	Name         string    `json:"name"`
// 	ModelName    string    `json:"model_name"`
// 	Instructions string    `json:"instructions"`
// 	FileID       string    `json:"file_id"`
// 	VectorID     string    `json:"vector_id"`
// 	LastEdit     time.Time `json:"last_edit"`
// }

// const aiTableSQL = `
// CREATE TABLE IF NOT EXISTS ai (
//     id SERIAL PRIMARY KEY,
//     name VARCHAR(255) UNIQUE NOT NULL,
// 	model_name VARCHAR(255) NOT NULL,
// 	instructions VARCHAR(255) NOT NULL,
// 	file_id VARCHAR(255),
//     vector_id VARCHAR(255),
// 	last_edit TIMESTAMP NOT NULL
// );`

// func create_ai_table() {
// 	conn := establish_connection()
// 	var err error
// 	defer conn.Close()
// 	defer log.Println("Conn Closed")

// 	// Create tables (if they don't exist)
// 	_, err = conn.Exec(context.Background(), aiTableSQL)
// 	if err != nil {
// 		log.Fatalf("Failed to create table: %v\n", err)
// 	}

// 	ai := AI{
// 		Name:         "default",
// 		ModelName:    Model_Name,
// 		Instructions: "You are to simulate being a customer and use the provided file to answer questions asked by the user who is playing the role of a waiter",
// 		LastEdit:     time.Now(),
// 	}

// 	create_ai(ai)
// }

// func create_ai(ai AI) error {
// 	conn := establish_connection()
// 	var err error
// 	defer conn.Close()
// 	defer log.Println("Conn Closed")

// 	name_ai, err := retrieve_ai_pass_conn(conn, ai.Name)
// 	if err != nil {
// 		log.Println("Valid name")
// 	}
// 	if name_ai.Name == ai.Name {
// 		err := errors.New("Name: " + ai.Name + " already in use\n")
// 		log.Println(err)
// 		return err
// 	}

// 	log.Println("Creating ai")

// 	var aiID int
// 	// Prepare the SQL statement
// 	insertSQL := `INSERT INTO ai (name, model_name, instructions, file_id, vector_id, last_edit)
//     	VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`

// 	// Execute the SQL statement using a prepared statement
// 	err = conn.QueryRow(context.Background(), insertSQL, ai.Name, ai.ModelName, ai.Instructions, ai.FileID, ai.VectorID, time.Now()).Scan(&aiID)
// 	if err != nil {
// 		log.Fatalf("Failed to insert data: %v\n", err)
// 		return err
// 	}

// 	log.Printf("AI created with ID: %d\n", aiID)
// 	return nil
// }

// // TODO: need to create a function that checks against personas and makes sure they are up to date
// func update_ai(ai AI) error {
// 	conn := establish_connection()
// 	var err error
// 	defer conn.Close()
// 	defer log.Println("Conn Closed")

// 	name_ai, err := retrieve_ai_pass_conn(conn, ai.Name)
// 	if err != nil {
// 		log.Println("AI with that name could not be found:", err)
// 		return err
// 	}

// 	// This checks to see if name is connected to a real ai
// 	if name_ai.Name != ai.Name {
// 		log.Printf("AI: %s does not exist\n", ai.Name)
// 		return err
// 	}
// 	// Prepare the SQL statement for updating the user's name
// 	updateAISQL := `UPDATE ai SET model_name = $1, instructions = $2, file_id = $3, vector_id = $4, last_edit = $5
// 		WHERE name = $6;`

// 	// Execute the SQL statement using a prepared statement
// 	_, err = conn.Exec(context.Background(), updateAISQL, ai.ModelName, ai.Instructions, ai.FileID, ai.VectorID, time.Now(), ai.Name)
// 	if err != nil {
// 		log.Fatalf("Failed to update AI's info: %v\n", err)
// 		return err
// 	}
// 	log.Println("Update Complete")
// 	return nil
// }

// // TODO: finish retrieve ai function
// func retrieve_ai(name string) (ai AI, err error) {
// 	conn := establish_connection()
// 	defer conn.Close()
// 	defer log.Println("Conn Closed")

// 	// Prepare the SQL statement for selecting the user's data
// 	selectUserSQL := `SELECT id, name, model_name, instructions, file_id, vector_id, last_edit
//     	FROM ai
//     	WHERE name = $1;`

// 	err = conn.QueryRow(context.Background(), selectUserSQL, name).Scan(
// 		&ai.AI_ID, //the id variable should not be used outside the backend
// 		&ai.Name,
// 		&ai.ModelName,
// 		&ai.Instructions,
// 		&ai.FileID,
// 		&ai.VectorID,
// 		&ai.LastEdit,
// 	)
// 	if err != nil {
// 		err = errors.New("Failed to retireve AI:" + err.Error())
// 		return
// 	}
// 	return
// }

// // TODO: finish retrieve ai function
// func retrieve_ai_pass_conn(conn *pgxpool.Pool, name string) (ai AI, err error) {
// 	// Prepare the SQL statement for selecting the user's data
// 	selectUserSQL := `SELECT id, name, model_name, instructions, file_id, vector_id, last_edit
//     	FROM ai
//     	WHERE name = $1;`

// 	err = conn.QueryRow(context.Background(), selectUserSQL, name).Scan(
// 		&ai.AI_ID, //the id variable should not be used outside the backend
// 		&ai.Name,
// 		&ai.ModelName,
// 		&ai.Instructions,
// 		&ai.FileID,
// 		&ai.VectorID,
// 		&ai.LastEdit,
// 	)
// 	if err != nil {
// 		err = errors.New("Failed to retireve AI:" + err.Error())
// 		return
// 	}
// 	return
// }

// func upload_menu(menu_name string, ai AI) error {
// 	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

// 	file, err := os.Open(filepath.Join(Menu_Path, menu_name))
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}
// 	defer file.Close()

// 	//Need to get old menu and delete it========================================================================
// 	//Or make a menu table so that the menu can be better catigorized also giving the ability for different menus

// 	// Get the file size
// 	stat, err := file.Stat()
// 	if err != nil {
// 		log.Println(err)
// 		return err
// 	}

// 	// Read the file into a byte slice
// 	bs := make([]byte, stat.Size())
// 	_, err = bufio.NewReader(file).Read(bs)
// 	if err != nil && err != io.EOF {
// 		log.Println("File not found: ", err)
// 		return err
// 	}
// 	log.Println(bs)

// 	if ai.FileID != "" {
// 		err = client.DeleteFile(context.Background(), ai.FileID)
// 		if err != nil {
// 			log.Println("Old file could not be deleted: ", err)
// 			return err
// 		}
// 	}

// 	if ai.VectorID != "" {
// 		_, err = client.DeleteVectorStore(context.Background(), ai.VectorID)
// 		if err != nil {
// 			log.Println("Old file could not be deleted: ", err)
// 			return err
// 		}
// 	}

// 	var file_upload openai.FileBytesRequest
// 	file_upload.Name = menu_name
// 	file_upload.Bytes = bs
// 	file_upload.Purpose = "assistants"

// 	file_return, err := client.CreateFileBytes(context.Background(), file_upload)
// 	if err != nil {
// 		log.Println("File not stored in openai: ", err)
// 		return err
// 	}

// 	ai.FileID = file_return.ID

// 	var vector_upload openai.VectorStoreRequest

// 	vector_upload.Name = "menu"
// 	vector_upload.FileIDs = append(vector_upload.FileIDs, file_return.ID)

// 	vector_return, err := client.CreateVectorStore(context.Background(), vector_upload)
// 	if err != nil {
// 		err = client.DeleteFile(context.Background(), ai.FileID)
// 		if err != nil {
// 			log.Println("Old file could not be deleted: ", err)
// 			return err
// 		}
// 		log.Println("File not transfered to vector: ", err)
// 		return err
// 	}

// 	ai.VectorID = vector_return.ID

// 	update_ai(ai)

// 	return nil
// }

// func create_menu_record(menu_data, menu_name string) (string, error) {
// 	file, err := os.OpenFile(filepath.Join(Menu_Path, menu_name+".json"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
// 	if err != nil {
// 		log.Println("Menu file not created: ", err)
// 		return "", err
// 	}
// 	defer file.Close()

// 	var jsonMap map[string]interface{}
// 	json.Unmarshal([]byte(menu_data), &jsonMap)

// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(jsonMap)
// 	if err != nil {
// 		log.Println("Data not encoded into file: ", err)
// 		err = os.Remove(filepath.Join(Menu_Path, menu_name+".json"))
// 		return "", err
// 	}

// 	return menu_name + ".json", nil
// }
