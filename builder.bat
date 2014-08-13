@echo off


echo RU-SUPPLIER & echo.
type LICENSE & echo.


echo Checking compiler...
where go
if %ERRORLEVEL% NEQ 0 (
	echo Error: compiler was not found
	goto :end
)


echo Making directory tree
if exist build (
	rd /q /s build
)
md build
md build\urls
md build\docs
copy docs build\docs


echo Compilation of ru-supplier...
copy eagle.ico build
copy LICENSE build
copy README.md build
copy rsrc.syso build
copy ru-supplier.manifest build
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check golang version and .go files
	goto :end
)
move ru-supplier.exe build


echo Compilation of url-generator...
copy urls\rsrc.syso build\urls
copy urls\urls.manifest build\urls
cd urls
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check url/*.go files
	goto :end
)
cd ..
move urls\urls.exe build\urls


echo RU-SUPPLIER successfully installed in: %CD%\build\
start build

:end