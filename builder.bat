
:: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: ::
::                                                       ::
::  Script for compilation of ru-supplier programm       ::
::  author: Pushkin Ivan <pushkin13@bk.ru>               ::
::  link: https://bitbucket.org/pushkin_ivan/ru-supplier ::
::                                                       ::
:: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: ::

@echo off

::
:: Printing Programm name, author, licence
::
echo RU-SUPPLIER & echo.
echo Author: Pushkin Ivan [pushkin13@bk.ru]
type LICENSE & echo.

::
:: Checking compiler
::
echo Checking compiler...
where go
if %ERRORLEVEL% NEQ 0 (
	echo Error: compiler was not found
	goto :end
)

::
:: Making directories and copy docs
::
echo Making directory tree
if exist build (
	rd /q /s build
)
md build
md build\urls
md build\docs
copy docs build\docs

::
:: Programm compilation and prog files coping
::
echo Compilation of ru-supplier...
copy eagle.ico build
copy LICENSE build
copy README.md build
copy rsrc.syso build
copy ru-supplier.manifest build
echo Please wait...
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check golang version and .go files
	goto :end
)
move ru-supplier.exe build

::
:: Compilation of url-generator
::
echo Compilation of url-generator...
copy urls\rsrc.syso build\urls
copy urls\urls.manifest build\urls
cd urls
echo Please wait...
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check url/*.go files
	goto :end
)
cd ..
move urls\urls.exe build\urls

::
:: Finalizing
::
echo RU-SUPPLIER successfully installed in: %CD%\build\
start build

:end
