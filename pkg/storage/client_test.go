package storage

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )
//
// func visit(path string, info os.FileInfo, err error) error {
// 	if err != nil {
// 		fmt.Println(err) // can't walk here,
// 		return nil       // but continue walking elsewhere
// 	}
//
// 	if info.IsDir() {
// 		fmt.Println("Directory:", path)
// 	} else {
// 		fmt.Println("File:", path)
// 	}
//
// 	return nil
// }
//
// func main() {
// 	root := "." // change this to the directory you want to start the recursive iteration from
//
// 	err := filepath.Walk(root, visit)
// 	if err != nil {
// 		fmt.Printf("error walking the path %v: %v\n", root, err)
// 	}
// }
