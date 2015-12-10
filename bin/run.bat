echo off

taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe

start "Lobby 1" E:\LiveServer\src\gameserver\gameserver.exe lobby 1
start "Auth 1" E:\LiveServer\src\gameserver\gameserver.exe auth 1
start "Master 1" E:\LiveServer\src\gameserver\gameserver.exe master 1
start "Gate 1" E:\LiveServer\src\gameserver\gameserver.exe gate 1
start "Connector 1" E:\LiveServer\src\gameserver\gameserver.exe connector 1
start "Game 1" E:\LiveServer\src\gameserver\gameserver.exe game 1
start "DataBase 1" E:\LiveServer\src\gameserver\gameserver.exe database 1
start "Society 1" E:\LiveServer\src\gameserver\gameserver.exe society 1