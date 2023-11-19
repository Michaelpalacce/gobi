package main

import "github.com/Michaelpalacce/gobi/pkg/database"

func main() {
    driver := database.Get()

    defer driver.Disconnect()
}
