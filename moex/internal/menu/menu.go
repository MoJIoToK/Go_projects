package menu

import (
	"bufio"
	"fmt"

	"github.com/olekukonko/tablewriter"
	"log/slog"
	"moex-app/internal/api"
	"moex-app/internal/service"
	"os"
	"strings"
)

func StartMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	log := slog.Default()
	log.Info("StartMenu")

	for {
		log.Info("Main menu")
		fmt.Println("\nВыберите действие:")
		fmt.Println("1 - Получить котировки по тикеру")
		fmt.Println("2 - Посмотреть список доступных инструментов")
		fmt.Println("3 - Выйти")

		fmt.Print("Введите номер действия: ")
		if !scanner.Scan() {
			slog.Error("Invalid input in main menu")
			fmt.Println("Пожалуйста повторите снова")
			continue
		}
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			slog.Info("User choices 1")
			fmt.Print("Введите тикер: ")
			if !scanner.Scan() {
				slog.Error("Invalid input in 1 menu")
				fmt.Println("Ошибка ввода")
				continue
			}
			ticker := scanner.Text()
			num := strings.Split(ticker, ",")

			jsonDoc, err := api.GetMarketData(ticker)
			if err != nil {
				fmt.Println(err)
			}

			resDoc, err := service.ExtractData(jsonDoc)
			if err != nil {
				fmt.Println(err)
			}

			printToConsole(resDoc, len(num))

		case "2":
			slog.Info("User choices 2")
			fmt.Println("Функция просмотра списка инструментов пока не реализована.")

		case "3":
			slog.Info("App exiting")
			fmt.Println("Выход из программы.")
			os.Exit(0)

		default:
			slog.Error("Invalid choice in main menu")
			fmt.Println("Please try again")
		}
	}
}

func printToConsole(doc []map[string]interface{}, num int) {

	log := slog.Default()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Наименование компании", "Тикер", "Последняя стоимость за акцию"})

	for i := 0; i < num; i++ {
		row := []string{
			doc[i]["SECNAME"].(string),
			doc[i]["SECID"].(string),
			fmt.Sprintf("%v", doc[i]["LAST"]),
		}
		table.Append(row)
	}

	table.SetAutoFormatHeaders(false)
	table.SetBorders(tablewriter.Border{Left: false, Top: false, Right: false, Bottom: false})

	log.Info("Print to console")
	table.Render()
}
