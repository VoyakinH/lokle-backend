package main

import (
	"context"
	// "database/sql"
	"encoding/json"
	// "fmt"
	// "github.com/xuri/excelize/v2"
	"net/http"
	"os"

	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

var ctx = context.Background()

const expCookieTime = 1382400

var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

// var psqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 	"185.225.35.60", 5432, "postgres", "1991055q", "glasha")

type Auth struct {
	Email    string
	Password string
}

type User struct {
	password string
	budgets  map[string]Budget
}

type Budget struct {
	balance           string
	last_payment_date string
	students          map[string]Student
}

type Student struct {
	lessons map[string]Lesson
}

type Lesson struct {
	date_time string
	teacher   string
	duration  string
	status    string
	note      string
}

// func importXLSX() {
// 	f, err := excelize.OpenFile("./bd.xlsm")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	db, _ := sql.Open("postgres", psqlconn)

// 	comp_id_cols, _ := f.GetCols("Сопоставление ID")
// 	budgets_cols, _ := f.GetCols("Budgets")
// 	payments_cols, _ := f.GetCols("Payments_Students")
// 	lessons_cols, _ := f.GetCols("Lessons")

// 	cic_students_index := 0
// 	cic_budgets_index := 1
// 	cic_logins_index := 4
// 	cic_passwords_index := 5

// 	bc_names_index := 0
// 	bc_balances_index := 1

// 	pc_dates_index := 1
// 	pc_types_index := 3
// 	pc_budgets_index := 5

// 	lc_date_index := 0
// 	lc_time_index := 2
// 	lc_student_index := 3
// 	lc_teacher_index := 4
// 	lc_duration_index := 5
// 	lc_status_index := 7
// 	lc_note_index := 13

// 	//insert_users := `INSERT INTO "user" ("login","password") VALUES `

// 	//n := 0
// 	balances_ref := map[string]string{}
// 	budget_dates_ref := map[string]string{}
// 	users := map[string]User{}

// 	for i, val := range budgets_cols[bc_balances_index] {
// 		if val != "" && i > 0 {
// 			balances_ref[budgets_cols[bc_names_index][i]] = val
// 		}
// 	}

// 	for i := len(payments_cols[pc_budgets_index]) - 1; i >= 1; i-- {
// 		budget_name := payments_cols[pc_budgets_index][i]
// 		budget_type := payments_cols[pc_types_index][i]
// 		budget_date := payments_cols[pc_dates_index][i]
// 		if budget_name != "" && budget_type == "Оплата блока" && budget_dates_ref[budget_name] == "" {
// 			budget_dates_ref[budget_name] = budget_date
// 		}
// 	}

// 	for i, val := range comp_id_cols[cic_logins_index] {
// 		if val != "" && i > 0 {
// 			login := comp_id_cols[cic_logins_index][i]
// 			password := comp_id_cols[cic_passwords_index][i]
// 			budget_name := comp_id_cols[cic_budgets_index][i]
// 			student_name := comp_id_cols[cic_students_index][i]
// 			if users[login].password == "" {
// 				users[login] = User{
// 					password: password,
// 					budgets:  map[string]Budget{},
// 				}
// 			}

// 			balance := balances_ref[budget_name]
// 			if balance == "" {
// 				balance = balances_ref[budget_name+" "]
// 			}

// 			if users[login].budgets[budget_name].balance == "" {
// 				users[login].budgets[budget_name] = Budget{
// 					balance:           balance,
// 					last_payment_date: budget_dates_ref[budget_name],
// 					students:          map[string]Student{},
// 				}
// 			}

// 			users[login].budgets[budget_name].students[student_name] = Student{
// 				lessons: map[string]Lesson{},
// 			}

// 			//users[comp_id_cols[comp_id_logins_index][i]] = User{
// 			//	login:    comp_id_cols[comp_id_logins_index][i],
// 			//	password: comp_id_cols[comp_id_passwords_index][i]}

// 			//if !strings.Contains(insert_users, comp_id_cols[comp_id_logins_index][i]) {
// 			//insert_users += `('` + comp_id_cols[comp_id_logins_index][i] + `','` +
// 			//	comp_id_cols[comp_id_passwords_index][i] + `'),`
// 			//}
// 			// comp_id_cols[comp_id_logins_index][n] = val
// 			// comp_id_cols[comp_id_passwords_index][n] = comp_id_cols[comp_id_passwords_index][i]
// 			// comp_id_cols[comp_id_budgets_index][n] = comp_id_cols[comp_id_budgets_index][i]
// 			// n++
// 		}
// 	}
// 	//insert_users = insert_users[:len(insert_users)-1] + `;`

// 	//comp_id_cols[comp_id_logins_index] = comp_id_cols[comp_id_logins_index][1:n]
// 	//comp_id_cols[comp_id_passwords_index] = comp_id_cols[comp_id_passwords_index][1:n]
// 	//comp_id_cols[comp_id_budgets_index] = comp_id_cols[comp_id_budgets_index][1:n]

// 	//db.Exec(`TRUNCATE TABLE "user" RESTART IDENTITY CASCADE;`)
// 	//db.Exec(insert_users)

// 	fmt.Print(users)
// 	// for _, col := range cols {
// 	// 	for _, rowCell := range col {
// 	// 		fmt.Print(rowCell, "\t")
// 	// 	}
// 	// }
// 	defer db.Close()
// 	f.Close()
// }

func createUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://185.225.34.197")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	var a Auth
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionID, _ := uuid.NewRandom()

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  sessionID.String(),
		MaxAge: expCookieTime,
	}

	_, err = rdb.SetNX(ctx, sessionID.String(), 0, expCookieTime*time.Second).Result()
	if err != nil {
		w.WriteHeader(523)
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			return
		}
		return
	}

	http.SetCookie(w, cookie)
}

func deleteUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://185.225.34.197")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	tokenCookie, err := r.Cookie("session-id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rdb.Del(ctx, tokenCookie.Value).Val()

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}

	http.SetCookie(w, cookie)
}

func checkUserSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://185.225.34.197")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	tokenCookie, err := r.Cookie("session-id")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	val := rdb.Exists(ctx, tokenCookie.Value).Val()
	if val < 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rdb.Expire(ctx, tokenCookie.Value, expCookieTime*time.Second)

	cookie := &http.Cookie{
		Name:   "session-id",
		Value:  tokenCookie.Value,
		MaxAge: expCookieTime,
	}

	http.SetCookie(w, cookie)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	mux := http.NewServeMux()

	// importXLSX()

	mux.HandleFunc("/login", createUserSession)
	mux.HandleFunc("/logout", deleteUserSession)
	mux.HandleFunc("/validate_user", checkUserSession)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		return
	}
}
