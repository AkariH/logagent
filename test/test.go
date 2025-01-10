package main

import (
    "fmt"
    "os"
    "strconv"

    "go.uber.org/zap"
)

func main() {
    // Check if the user provided the skip count as an argument
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run test.go <n>")
        return
    }

    // Parse the skip count from command-line arguments
    n, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println("Please provide a valid integer for n.")
        return
    }

    // Loop from 0 to n
    for i := 0; i <= n; i++ {
        // Create a logger with AddCallerSkip(i)
        logger, err := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(i))
        if err != nil {
            fmt.Printf("Failed to create logger with AddCallerSkip(%d): %v\n", i, err)
            continue
        }
        defer logger.Sync()

        // Log a message
        logger.Info(fmt.Sprintf("This is a log message with AddCallerSkip(%d)", i))
    }
}