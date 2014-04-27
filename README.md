# mdbToCsv takes in a .mdb file and converts the tables within it to .csv files


# Usage
First, make sure you have the MS Access Database Engine installed: http://www.microsoft.com/en-us/download/details.aspx?id=13255

```
go install github.com/oztu/mdbToCsv
mdbToCsv myMdbFile.mdb tableName1 tableName2 tableName3 
```

This will create tableName1.csv, tableName2.csv, and tableName3.csv in the current directory.