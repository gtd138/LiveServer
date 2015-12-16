echo off

taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe
taskkill /f /im gameserver.exe

start "Lobby 1" C:\Github\LiveServer\bin\gameserver.exe lobby 1
start "Auth 1" C:\Github\LiveServer\bin\gameserver.exe auth 1
start "Master 1" C:\Github\LiveServer\bin\gameserver.exe master 1
start "Gate 1" C:\Github\LiveServer\bin\gameserver.exe gate 1
start "Connector 1" C:\Github\LiveServer\bin\gameserver.exe connector 1
start "Game 1" C:\Github\LiveServer\bin\gameserver.exe game 1
start "DataBase 1" C:\Github\LiveServer\bin\gameserver.exe database 1
start "Society 1" C:\Github\LiveServer\bin\gameserver.exe society 1