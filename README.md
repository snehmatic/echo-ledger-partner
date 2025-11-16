# Echo Ledger - Partner

Echo Ledger - Partner is a personal tool to automate transation/expense tracking. It uses exported expense report (CSV) from tools like Axio, and filters and presents higher level info.

Currently Partner is a basic `Go` script but in future can (and planned to) be turned into a bridge between low level trackers (Axio, Excel, etc.) and high level trackers (Net worth, Annual stats).

## Pre-req: 

1. `go` (1.24.x)
2. exported `.csv` file with data (currently only Axio is supported, might not work with other tools) 

## Run local

`go run .`
