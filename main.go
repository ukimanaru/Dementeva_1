package main
import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Expense struct{
	ID		int
	Amount		float64
	Description		string
	Date	string
}

var expenses=[]Expense{
	{ID: 1, Amount: 250, Description: "Продукты", Date: "2026-09-20"},
	{ID: 2, Amount: 120, Description: "Транспорт", Date: "2026-09-21"},
	{ID: 3, Amount: 500, Description: "Покупки", Date: "2026-09-22"},
}

var templates=template.Must(template.ParseFiles(
	"templates/layout.html",
	"templates/expenses.html",
))

var nextID=4

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method!=http.MethodGet{
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
http.Redirect(w,r, "/expenses", http.StatusSeeOther)
}

func aboutPage(w http.ResponseWriter, r *http.Request) {
	if r.Method !=http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("Мой личный трекер расходов"))
}

func pingPage(w http.ResponseWriter, r *http.Request) {
	if r.Method !=http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Write([]byte("pong"))
}

func expensesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method==http.MethodGet{
		err:=templates.ExecuteTemplate(w, "layout.html", expenses)
		if err!=nil {
			http.Error(w,"Ошибка отображения страницы", http.StatusInternalServerError)
			return
		}
		return
	}

	if r.Method==http.MethodPost{
		amountText := r.FormValue("amount")
		description:=r.FormValue("description")
		date:=r.FormValue("date")

		amount,err:=strconv.ParseFloat(amountText,64)

		if err!=nil || amount<=0 || description== "" || date== "" {
			http.Error(w,"Проверьте введённые данные", http.StatusBadRequest)
			return
		}

		newExpense:=Expense{
			ID:		len(expenses)+1,
			Amount:		amount,
			Description:		description,
			Date:		date,
		}

		expenses=append(expenses, newExpense)

		http.Redirect(w,r,"/expenses", http.StatusSeeOther)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/ping", pingPage)
	http.HandleFunc("/expenses", expensesPage)

	log.Println("Сервер запущен: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}