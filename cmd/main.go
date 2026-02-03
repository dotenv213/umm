package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dotenv213/umm/internal/userstore"
)

func readLine(scanner *bufio.Scanner, prompt string) string {
	fmt.Print(prompt)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func main() {
	store, err := userstore.NewDb("users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	scanner := bufio.NewScanner(os.Stdin)
	ctx := context.Background()

	for {
		fmt.Println("\n--- User Management System ---")
		fmt.Println("1. Create User")
		fmt.Println("2. List All Users")
		fmt.Println("3. Update User")
		fmt.Println("4. Delete User")
		fmt.Println("5. Exit")
		fmt.Println("Select an option: ")

		scanner.Scan()
		choice := scanner.Text()
		switch choice {
		case "1":
			uname := readLine(scanner, "Enter Username: ")
			email := readLine(scanner, "Enter Email: ")
			if uname == "" || email == "" {
				fmt.Println("username and email are required")
				continue
			}
			u := &userstore.User{Username: uname, Email: email}
			if err := store.Create(ctx, u); err != nil {
				fmt.Printf("Error %v\n", err)
			} else {
				fmt.Println("User Created!")
			}
		case "2":
			users, err := store.ListAll(ctx)
			if err != nil {
				fmt.Println("failed to list users:", err)
				continue
			}
			fmt.Println("\n  ID  |  Username  |  Email  | Created at  ")
			for _, u := range users {
				fmt.Printf("%-3d  |  %-10s  |  %s  |  %v  \n", u.ID, u.Username, u.Email, u.CreatedAt)
			}

		case "3":
			idStr := readLine(scanner, "Enter user ID: ")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				fmt.Println("invalid id")
				continue
			}

			u, err := store.GetById(ctx, id)
			if err != nil {
				fmt.Println("User not found")
				continue
			}
			newU := readLine(scanner, fmt.Sprintf("Username [%s]: ", u.Username))
			if newU != "" {
				u.Username = newU
			}

			newE := readLine(scanner, fmt.Sprintf("Email [%s]: ", u.Email))
			if newE != "" {
				u.Email = newE
			}

			if err := store.Update(ctx, u); err != nil {
				fmt.Printf("Update failed: %v\n", err)
			} else {
				fmt.Println("Updated successfully!")
			}
		case "4":
			idStr := readLine(scanner, "Enter a user ID to delete: ")
			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				fmt.Println("Invalid ID format")
				continue
			}

			u, err := store.GetById(ctx, id)
			if err != nil {
				fmt.Println("User not found")
				continue
			}

			confirm := readLine(scanner, "Are you sure you want to delete? (y/n): ")
			if confirm != "y" {
				continue
			} else {
				if err := store.Delete(ctx, u.ID); err != nil {
					fmt.Println("Delete failed")
				} else {
					fmt.Println("User deleted successfuly")
				}
			}
		case "5":
			fmt.Println("Exiting program...")
			return
		}
	}
}
