package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-kivik/couchdb/v3"
	kivik "github.com/go-kivik/kivik/v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type StudentModel struct {
	ID          string `json:"_id"`
	StuId       string `json:"stu_id"`
	Name        string `json:"name"`
	Sex         string `json:"sex"`
	Grade       string `json:"grade"`
	PhoneNumber string `json:"phone_number"`
}

type StudentModelRequest struct {
	StuId       string `json:"stu_id"`
	Name        string `json:"name"`
	Sex         string `json:"sex"`
	Grade       string `json:"grade"`
	PhoneNumber string `json:"phone_number"`
}
type StudentModelImage struct {
	StuId       string `json:"stu_id"`
	Name        string `json:"name"`
	Sex         string `json:"sex"`
	Grade       string `json:"grade"`
	PhoneNumber string `json:"phone_number"`
	Attachments string `json:"attachments"`
}
type lener interface {
	Len() int
}
type stater interface {
	Stat() (os.FileInfo, error)
}

func readerSize(in io.Reader) (int64, io.Reader, error) {
	if ln, ok := in.(lener); ok {
		return int64(ln.Len()), in, nil
	}
	if st, ok := in.(stater); ok {
		info, err := st.Stat()
		if err != nil {
			return 0, nil, err
		}
		return info.Size(), in, nil
	}
	content, err := ioutil.ReadAll(in)
	if err != nil {
		return 0, nil, err
	}
	buf := bytes.NewBuffer(content)
	return int64(buf.Len()), ioutil.NopCloser(buf), nil
}

func InsertNewDocs(ctx *gin.Context) {
	var body StudentModel
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Input"})
		return
	}
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")

	doc := map[string]interface{}{
		"_id":          body.ID,
		"stu_id":       body.StuId,
		"name":         body.Name,
		"sex":          body.Sex,
		"grade":        body.Grade,
		"phone_number": body.PhoneNumber,
	}
	fmt.Println(doc)
	rev, err := db.Put(context.TODO(), body.ID, doc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Document inserted successfully", "revision": rev})
}
func GetAllDocs(ctx *gin.Context) {
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")
	rows, err := db.AllDocs(context.TODO(), kivik.Options{"include_docs": true})
	var docs []StudentModelRequest
	for rows.Next() {
		var row struct {
			ID  string              `json:"id"`
			Doc StudentModelRequest `json:"doc"`
		}
		if err := rows.ScanDoc(&row.Doc); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
			return
		}
		docs = append(docs, row.Doc)
	}
	ctx.JSON(http.StatusOK, docs)
}
func GetDocsById(ctx *gin.Context) {
	id := ctx.Param("id")
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")

	doc := db.Get(context.TODO(), id)
	var student StudentModelRequest
	fmt.Println(student)
	if err := doc.ScanDoc(&student); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, student)
}
func UpdateDoc(ctx *gin.Context) {
	var body StudentModel
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")

	currentDoc := db.Get(context.Background(), id)

	// Add the current revision to the update
	var docMeta map[string]interface{}
	if err := currentDoc.ScanDoc(&docMeta); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
		return
	}
	fmt.Println(docMeta)
	//doc["_rev"] = docMeta["_rev"]
	doc := map[string]interface{}{
		"_id":          body.ID,
		"stu_id":       body.StuId,
		"name":         body.Name,
		"sex":          body.Sex,
		"grade":        body.Grade,
		"phone_number": body.PhoneNumber,
		"_rev":         docMeta["_rev"],
	}
	//fmt.Println(doc)
	rev, err := db.Put(context.Background(), body.ID, doc)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document" + rev})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Document updated successfully", "revision": rev})
}
func DeleteDoc(ctx *gin.Context) {
	var body StudentModel
	id := ctx.Param("id")
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.Background(), "student_db")
	currentDoc := db.Get(context.Background(), id)
	var docMeta map[string]interface{}
	if err := currentDoc.ScanDoc(&docMeta); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
		return
	}

	doc := map[string]interface{}{
		"_id":      id,
		"_rev":     docMeta["_rev"],
		"_deleted": true,
	}
	fmt.Println(doc)
	rev, err := db.Put(context.Background(), body.ID, doc)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update document" + rev})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Document updated successfully", "revision": rev})

}
func NewAttachment(filename, contentType string, content io.Reader, size ...int64) (*kivik.Attachment, error) {
	var filesize int64
	if len(size) > 0 {
		filesize = size[0]
	} else {
		var err error
		filesize, content, err = readerSize(content)
		if err != nil {
			return nil, err
		}
	}
	rc, ok := content.(io.ReadCloser)
	if !ok {
		rc = ioutil.NopCloser(content)
	}
	return &kivik.Attachment{
		Filename:    filename,
		ContentType: contentType,
		Content:     rc,
		Size:        filesize,
	}, nil
}
func UploadImage(ctx *gin.Context) {
	id := ctx.Param("id")

	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")

	// Get the current document
	currentDoc := db.Get(context.Background(), id)
	var docMeta map[string]interface{}
	if err := currentDoc.ScanDoc(&docMeta); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
		return
	}

	// Retrieve the file from the request
	file, _, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file from request: " + err.Error()})
		return
	}
	defer file.Close()

	// Create the attachment
	attachment, err := NewAttachment("photo.jpg", "image/jpeg", file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attachment: " + err.Error()})
		return
	}
	atts := map[string]*kivik.Attachment{
		"photo.jpg": attachment,
	}

	// Prepare the updated document
	updatedDoc := map[string]interface{}{
		"_id":          id,
		"_rev":         docMeta["_rev"],
		"_attachments": atts,
	}
	// Merge the existing document fields with the new attachment
	for key, value := range docMeta {
		if key != "_id" && key != "_rev" && key != "_attachments" {
			updatedDoc[key] = value
		}
	}

	// Save the updated document
	rev, err := db.Put(context.Background(), id, updatedDoc)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert document"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Document updated successfully", "revision": rev})
}
func GetImage(ctx *gin.Context) {
	id := ctx.Param("id")
	atts := ctx.Param("atts")
	fmt.Println("ID:", id, "Attachment:", atts)
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	db := client.DB(context.TODO(), "student_db")

	attachment, err := db.GetAttachment(context.TODO(), id, atts)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found: " + err.Error()})
		return
	}
	defer attachment.Content.Close()

	// Read the attachment into a byte slice
	data, err := ioutil.ReadAll(attachment.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read attachment: " + err.Error()})
		return
	}

	// Set the appropriate content type
	ctx.Writer.Header().Set("Content-Type", attachment.ContentType)
	ctx.Writer.WriteHeader(http.StatusOK)

	// Write the data to the response
	ctx.Writer.Write(data)
}
func Filter(ctx *gin.Context) {
	// Initialize CouchDB client
	client, err := kivik.New("couch", "http://admin:123@localhost:5985/")
	if err != nil {
		log.Fatalf("Failed to connect to CouchDB: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to CouchDB"})
		return
	}
	//defer client.Close()

	// Access the database
	db := client.DB(context.TODO(), "student_db")

	// Define filter criteria (you can modify this as needed)
	name := ctx.Query("name")
	phone_number := ctx.Query("phone_number")

	// If both keys are required to filter, ensure they're provided
	if name == "" || phone_number == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name and phone_number parameters are required"})
		return
	}

	// Create query options
	opts := map[string]interface{}{
		"key":          []interface{}{name, phone_number}, // Composite key
		"include_docs": true,
	}

	// Perform the query
	rows, err := db.Query(context.TODO(), "_design/my_design", "by_name_phone", opts)
	if err != nil {
		log.Fatalf("Query error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Query error: " + err.Error()})
		return
	}
	defer rows.Close() // Ensure rows are closed when done

	// Collect results
	var docs []StudentModelRequest

	for rows.Next() {
		var row struct {
			ID  string              `json:"id"`
			Doc StudentModelRequest `json:"doc"`
		}
		if err := rows.ScanDoc(&row.Doc); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse document: " + err.Error()})
			return
		}
		docs = append(docs, row.Doc)
	}

	if err := rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Rows error: " + err.Error()})
		return
	}

	// Return filtered results
	ctx.JSON(http.StatusOK, docs)
	//http://localhost:8080/student_db/filter?name=meng&phone_number=0938271
}

func main() {
	r := gin.Default()
	r.POST("/student_db", InsertNewDocs)
	r.GET("/student_db", GetAllDocs)
	r.GET("/student_db/:id", GetDocsById)
	r.PUT("/student_db/update/:id", UpdateDoc)
	r.PUT("/student_db/delete/:id", DeleteDoc)
	r.PUT("/student_db/upload/:id", UploadImage)
	r.GET("/student_db/image/:id/:atts", GetImage)
	r.GET("/student_db/filter", Filter)

	r.Run(":8080")
}
