# pgsummary
This project builds executables that can be used to do the following:
 - write a summary report to a json file
 - compare two summary json files

NOTE:  The tool expects to get database credentials form AWS Secret Manager

The tool was build with AWS RDS postgres in mind, but probably works pretty broadly because it's simple.  


## Installation

Download the tarball and unpack it. You should get something like
```
/
├── darwin
│ └── amd64
│     ├── pgcompare
│     └── pgreport
├── linux
│ └── amd64
│     ├── pgcompare
│     └── pgreport
├── pgsummary-0.0.0.tar.gz
└── version.txt
```
## Usage

Create a summary 
```shell
pgreport \
-host=my.database.fqdn \
-dbname=postgres \
-secretId=my/secret/id/path \
# These are the keys in the secret JSON document
# we extract the values and use those
-secretUsernameKey=master_username \
-secretPasswordKey=master_password
```