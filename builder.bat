
:: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: ::
::                                                       ::
::  Script for compilation of ru-supplier programm       ::
::  Author: Pushkin Ivan <pushkin13@bk.ru>               ::
::  Link: https://bitbucket.org/pushkin_ivan/ru-supplier ::
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
where /Q go
if %ERRORLEVEL% NEQ 0 (
	echo Error: compiler was not found
	echo Please install Golang
	goto :end
)
if "%GOPATH%" == "" (
	echo Please define GOPATH variable
	goto :end
)

::
:: Checking environment
::
echo Checking environment...
where /Q hg
if %ERRORLEVEL% NEQ 0 (
	echo Error: mercurial was not found
	echo Please install mercurial
	goto :end
)
where /Q git
if %ERRORLEVEL% NEQ 0 (
	echo Error: git was not found
	echo Please install git and create git in PATH
	goto :end
)

::
:: Checking packages
::
echo Checking packages
if not exist "%GOPATH%\src\code.google.com/p/go-charset/charset" (
	echo Golang package code.google.com/p/go-charset/charset was not found
	echo Downloading package code.google.com/p/go-charset/charset...
	go get code.google.com/p/go-charset/charset
)
echo Package code.google.com/p/go-charset/charset... exists

if not exist "%GOPATH%\src\code.google.com/p/go-charset/data" (
	echo Golang package code.google.com/p/go-charset/data was not found
	echo Downloading package code.google.com/p/go-charset/data...
	go get code.google.com/p/go-charset/data
)
echo Package code.google.com/p/go-charset/data... exists

if not exist "%GOPATH%\src\github.com/gorilla/feeds" (
	echo Golang package github.com/gorilla/feeds was not found
	echo Downloading package github.com/gorilla/feeds...
	go get github.com/gorilla/feeds
)
echo Package github.com/gorilla/feeds... exists

if not exist "%GOPATH%\src\github.com/lxn/walk" (
	echo Golang package github.com/lxn/walk was not found
	echo Downloading package github.com/lxn/walk...
	go get github.com/lxn/walk
)
echo Package github.com/lxn/walk... exists

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

::
:: End builder.bat file
::
:end
