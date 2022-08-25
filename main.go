package main

/*
	1- Para instalar una caracteristicas necesarias para go
	usamos el comando </ go mod init directorioRaiz (app) >

	2- Para hacer una conexion con la base de datos nos descargamos un driver
	que esta en github con el comando </ go get -u github.com/go-sql-driver/mysql >
*/
/*
	"fmt": es para mostrar datos formateados
	"log": mostrar mensajes en la consola
 	"net/http": para manejar las peticiones http
	"html/template": mostrar las plantilla html com la extencion tmpl
*/

import (
	// "fmt"
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

/*
	Creramos la conexion a la base de datos
*/

func connection() (connection *sql.DB) {
	Driver := "mysql"
	User := "root"
	Pass := ""
	NameDB := "congreso"

	connection, err := sql.Open(Driver, User+":"+Pass+"@tcp(127.0.0.1)/"+NameDB)

	if err != nil {
		panic(err.Error())
	}

	return connection
}

/*
	Esta variable lee las plantillas o vistas
*/

var views = template.Must(template.ParseGlob("views/*"))

/*
	Esta funcion recibe la configuracion de
	del servidor
*/
func main() {
	// directorio raiz /directorio y Function
	http.HandleFunc("/", Home)
	http.HandleFunc("/tableRecord", TableRecord)
	http.HandleFunc("/record", Record)

	http.HandleFunc("/insertRecordPayment", InsertRecordPayment)

	http.HandleFunc("/deleteRecordPayment", DeleteRecordPayment)

	http.HandleFunc("/editRecordPayment", EditRecordPayment)

	http.HandleFunc("/updateRecordPayment", UpdateRecordPayment)
	// imprime en una linea
	log.Println("Servidor corriendo")
	// escucha el servidor en un puerto
	http.ListenAndServe(":8080", nil)
}

/* creamos un prototipo o estrutura de los registros y pagos */
type RecordPayment struct {
	Id, Idrecord, Idpayment int
}

/* update */
func UpdateRecordPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// recibimos los argumentos
		id := r.FormValue("id")
		record := r.FormValue("id_record")
		payment := r.FormValue("id_payment")

		// iniciamos la conexion
		connectActive := connection()
		updateRecordPayment, err := connectActive.Prepare(" UPDATE record_payment SET id_record=?, id_payment=? WHERE id=? ")

		if err != nil {
			panic(err.Error())
		}

		updateRecordPayment.Exec(record, payment, id)

		http.Redirect(w, r, "/tableRecord", 301)
	}
}

/*Form de edit un registro*/
func EditRecordPayment(w http.ResponseWriter, r *http.Request) {
	idRecordPayment := r.URL.Query().Get("id")

	// iniciamos la conexion
	connectActive := connection()
	oneRecordPayment, err := connectActive.Query(" SELECT * FROM record_payment WHERE id=?", idRecordPayment)

	if err != nil {
		panic(err.Error())
	}
	// esturctura
	recordPayment := RecordPayment{}

	for oneRecordPayment.Next() {
		var id, idrecord, idpayment int

		// capturamos un error
		err = oneRecordPayment.Scan(&id, &idrecord, &idpayment)
		if err != nil {
			panic(err.Error())
		}
		// si no hay error asignamos los datos a la estructura
		recordPayment.Id = id
		recordPayment.Idrecord = idrecord
		recordPayment.Idpayment = idpayment

	}
	// fmt.Println(recordPayment)
	// http.Redirect(w, r, "/tableRecord", 301)

	views.ExecuteTemplate(w, "editRecord", recordPayment)
}

/*Eliminamos un registro*/
func DeleteRecordPayment(w http.ResponseWriter, r *http.Request) {
	idRecordPayment := r.URL.Query().Get("id")
	// fmt.Println(idRecordPayment)

	// iniciamos la conexion
	connectActive := connection()
	insertRecordPayment, err := connectActive.Prepare(" DELETE FROM record_payment WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	insertRecordPayment.Exec(idRecordPayment)

	http.Redirect(w, r, "/tableRecord", 301)
}

/*
	Esta es una funciona master para
	recibir las peticiones de los clientes
*/
func Home(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hola mundo de Golang")
	views.ExecuteTemplate(w, "home", nil)
}

func TableRecord(w http.ResponseWriter, r *http.Request) {

	connectionActive := connection()
	queryRecordPayment, err := connectionActive.Query(" SELECT * FROM record_payment ")

	if err != nil {
		panic(err.Error())
	}

	/*asignamos la estructura*/
	recordPayment := RecordPayment{}
	arrayRecordPayment := []RecordPayment{}

	for queryRecordPayment.Next() {
		var id, idrecord, idpayment int

		// capturamos un error
		err = queryRecordPayment.Scan(&id, &idrecord, &idpayment)
		if err != nil {
			panic(err.Error())
		}
		// si no hay error asignamos los datos a la estructura
		recordPayment.Id = id
		recordPayment.Idrecord = idrecord
		recordPayment.Idpayment = idpayment

		// despues agregamos la estructura al arreglo
		arrayRecordPayment = append(arrayRecordPayment, recordPayment)

	}

	// fmt.Println(arrayRecordPayment)

	views.ExecuteTemplate(w, "tableRecord", arrayRecordPayment)
}

func Record(w http.ResponseWriter, r *http.Request) {

	/* realizamos una insercion de prueba*/
	// connectActive:= connection ()
	// insertRecord, err:= connectActive.Prepare(" INSERT INTO record_payment(id_record, id_payment) VALUE(2, 2) ")

	// if err != nil {
	// 	panic(err.Error())
	// }

	// insertRecord.Exec()

	views.ExecuteTemplate(w, "record", nil)
}

func InsertRecordPayment(w http.ResponseWriter, r *http.Request) {

	if r != nil {
		// recibimos los argumentos
		record := r.FormValue("id_record")
		payment := r.FormValue("id_payment")

		// iniciamos la conexion
		connectActive := connection()
		insertRecordPayment, err := connectActive.Prepare(" INSERT INTO record_payment(id_record, id_payment) VALUE(?,?) ")

		if err != nil {
			panic(err.Error())
		}

		insertRecordPayment.Exec(record, payment)

		http.Redirect(w, r, "/", 301)
	}

}
