package main

import (
	"fmt"
	"os"

	"github.com/SecDbg/memzip"
)

func main() {
	archive := memzip.NewZipArchive()

	// Add existing file or directory with an optional destination path
	err := archive.AddPath("your/file/path/here", "path/in/archive")
	if err != nil {
		fmt.Println("Error adding file or directory:", err)
		return
	}

	// Create a file directly
	err = archive.CreateFile("Hello.txt", "Hello! My name is John Doe.")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	zipBytes, err := archive.Close()
	if err != nil {
		fmt.Println("Error closing ZIP archive:", err)
		return
	}

	// Optionally, save the ZIP archive to a file
	err = os.WriteFile("output.zip", zipBytes, 0644)
	if err != nil {
		fmt.Println("Error writing ZIP file:", err)
		return
	}

	fmt.Println("ZIP archive created successfully.")
}
