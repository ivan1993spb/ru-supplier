
:: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: ::
::                                                       ::
::  Script for compilation of ru-supplier program        ::
::  Author: Pushkin Ivan <pushkin13@bk.ru>               ::
::  Link: https://github.com/ivan1993spb/ru-supplier     ::
::                                                       ::
:: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: :: ::

@echo off

::
:: Printing programm name, author, licence agreement
::
echo RU-SUPPLIER & echo.
echo Author: Pushkin Ivan (pushkin13@bk.ru) & echo.
type LICENSE & echo. & echo.

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
echo Compiler... ok

::
:: Checking environment
::
echo Checking environment...
if "%GOPATH%" == "" (
	echo Please define GOPATH env variable
	goto :end
)
echo GOPATH: (%GOPATH%) ... ok
where /Q hg
if %ERRORLEVEL% NEQ 0 (
	echo Error: mercurial was not found
	echo Please install mercurial
	goto :end
)
echo mercurial... ok
where /Q git
if %ERRORLEVEL% NEQ 0 (
	echo Error: git was not found
	echo Please install git command line tool
	goto :end
)
echo git... ok

::
:: Checking source directories
::
echo Checking source directories...
if not exist go-source\ru-supplier (
	echo not found source folder go-source\ru-supplier
	goto :end
)
if not exist go-source\urls (
	echo not found source folder go-source\urls
	goto :end
)
if not exist docs (
	echo not found folder docs
	goto :end
)
if not exist src (
	echo not found folder src
	goto :end
)
echo Source directories... ok

::
:: Checking packages
::
echo Checking packages...

if not exist "%GOPATH%\src\code.google.com/p/go-charset/charset" (
	echo Downloading package code.google.com/p/go-charset/charset...
	go get code.google.com/p/go-charset/charset
)
echo Package code.google.com/p/go-charset/charset... ok

if not exist "%GOPATH%\src\code.google.com/p/go-charset/data" (
	echo Downloading package code.google.com/p/go-charset/data...
	go get code.google.com/p/go-charset/data
)
echo Package code.google.com/p/go-charset/data... ok

if not exist "%GOPATH%\src\github.com/gorilla/feeds" (
	echo Downloading package github.com/gorilla/feeds...
	go get github.com/gorilla/feeds
)
echo Package github.com/gorilla/feeds... ok

if not exist "%GOPATH%\src\github.com/lxn/walk" (
	echo Downloading package github.com/lxn/walk...
	go get github.com/lxn/walk
)
echo Package github.com/lxn/walk... ok

echo.

:: inform that that's ok
echo ---------------------------------------
echo Your system is ready to compile program
echo ---------------------------------------& echo.

::
:: Compilation of ru-supplier
::
echo Compilation of ru-supplier...
copy src\common.manifest go-source\ru-supplier\ru-supplier.manifest
copy src\rsrc.syso go-source\ru-supplier
cd go-source\ru-supplier
echo Please wait...
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check golang version and .go files
	goto :end
)
echo Success!& echo.

del ru-supplier.manifest
del rsrc.syso
move ru-supplier.exe ..\..
cd ..\..

::
:: Compilation of url-generator
::
echo Compilation of url-generator...
copy src\common.manifest go-source\urls\urls.manifest
copy src\rsrc.syso go-source\urls
cd go-source\urls
echo Please wait...
go build -ldflags="-H windowsgui"
if %ERRORLEVEL% NEQ 0 (
	echo Compilation error: check golang version and .go files
	goto :end
)
echo Success!& echo.

del urls.manifest
del rsrc.syso
move urls.exe ..\..
cd ..\..

::
:: Finalizing
::
echo.
echo ---------------------------------
echo RU-SUPPLIER successfully compiled
echo ---------------------------------
echo.

::
:: End builder.bat file
::
:end

pause
