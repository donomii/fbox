del *.exe
del *.exe~
go build ..\vort\vort.go
"c:\Program Files (x86)\Inno Script Studio\ISStudio.exe" -compile installer.iss
echo Done!
