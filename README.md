# linkchain-explorer

## Mysql

Download [mysql](https://www.mysql.com/downloads/) if you haven't already

Config your mysql dsn at file `db/db.go`
```
...
var dsn = "dbusername:dbpassword@tcp(dbhost:dbport)/dbname?charset=utf8&parseTime=true"
...
```

Import tables `blocks` `transactions` `tickets`

The sql file in the directory `db`

## Installation (Linux or Mac)

Build program

`./build.sh`

Run

`lcexplorer`

Run nohup

`nohup lcexplorer >/**/*/debug.log 2>&1 &`

Then the server running at `http://127.0.0.1:9100`

## Installation (Windows)

Build program

`./build.bat`

Run (Windows)

`lcexplorer.exe`

Then the server running at `http://127.0.0.1:9100`