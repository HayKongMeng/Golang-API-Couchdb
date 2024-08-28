# CouchDB Student API

This Go application provides a RESTful API for managing student documents in a CouchDB database. It allows you to perform CRUD operations on student records, including inserting, updating, deleting, and retrieving documents. Additionally, it supports uploading and retrieving image attachments associated with student records.

## Features

- **Insert New Document**: Add a new student record to the database.
- **Get All Documents**: Retrieve all student records from the database.
- **Get Document by ID**: Fetch a specific student record by its ID.
- **Update Document**: Modify an existing student record.
- **Delete Document**: Remove a student record from the database.
- **Upload Image**: Attach an image to a student record.
- **Get Image**: Retrieve an image attached to a student record.
- **Filter Documents**: Query student records based on name and phone number.

## Prerequisites

- Go 1.18 or higher
- CouchDB (make sure it's installed and running)
- Gin-Gonic framework
- Go-Kivik CouchDB driver

## Setup

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd <repository-directory>

## API Endpoints
1. **Insert New Document**

   - Endpoint: POST /student_db
   - Description: Insert a new student record.
   - Request Body:
       ```
     {
     "_id": "stu-8",
     "stu_id": "stu-8",
     "name": "meng",
     "sex": "male",
     "grade": "a3",
     "phone_number": "0938271"
     }

2. Get All Documents
   - Endpoint: GET /student_db
   - Description: Retrieve all student records.

3. Get Document by ID

   - Endpoint: GET /student_db/:id
   - Description: Fetch a student record by its ID. 
   - URL Parameters:
       - id: The ID of the student document.
4. Update Document
   - Endpoint: PUT /student_db/update/:id
   - Description: Update an existing student record.

### URL Parameters:
   - id: The ID of the student document to update.
   - 
   Request Body:

    ```
    {
      "stu_id": "stu-8",
      "name": "meng",
      "sex": "male",
      "grade": "a3",
      "phone_number": "0938271"
    }

5. Delete Document

   - Endpoint: PUT /student_db/delete/:id
   - Description: Delete a student record.
   - URL Parameters:
     - id: The ID of the student document to delete.

Request Body:

    ```
    {
      "stu_id": "stu-8"
    }

6. Upload Image
   - Endpoint: PUT /student_db/upload/:id
   - Description: Upload an image and attach it to a student record.
   - URL Parameters:
       - id: The ID of the student document to which the image will be attached.
   - Form Data:
       - file: The image file to upload.

7. Get Image
   - Endpoint: GET /student_db/image/:id/:atts
   - Description: Retrieve an image attached to a student record.
   - URL Parameters:
     - id: The ID of the student document.
     - atts: The name of the attachment (e.g., photo.jpg).

8. Filter Documents
   - Endpoint: GET /student_db/filter
   - Description: Filter student records by name and phone number.
   - Query Parameters:
     - name: The name of the student to filter by.
     - phone_number: The phone number to filter by.

# Error Handling
    400 Bad Request: The request body or parameters are invalid.
    404 Not Found: The requested document or attachment does not exist.
    500 Internal Server Error: An unexpected error occurred.

# License
### This project is licensed under the MENG License.

# How to Use:
- **Clone the Repository**: Replace `<repository-url>` with the URL of your repository.
- **Run the Application**: Ensure CouchDB is properly configured and running.
- **Test the Endpoints**: Use cURL or Postman to interact with the API.

Feel free to modify and expand this README according to your needs and any additional features or setup instructions you may have.
