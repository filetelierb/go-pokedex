package main

import (
	"bufio"
	"os"
	"fmt"
)




func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		
		if scanner.Scan(){
			userInput := scanner.Text()
			splitUserInput := cleanInput(userInput)
			commandString := splitUserInput[0]
			parameters := splitUserInput[1:]
			if val,ok := cliCommands[commandString]; ok {
				val.callback(&initialConfigs, parameters)
			}
			
			//fmt.Printf("Your command was: %s\n",strings.ToLower(cleanInput(scanner.Text())[0]))
			
			
		} else {
			break
		}
		
	}

}
