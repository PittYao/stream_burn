chcp 65001
@echo off 

set _port=8084

set _task=stream_burn.exe
set _des=startup.bat

echo.
echo ========================================
echo == 查询视频刻录服务的状态==
echo == 每间隔10秒进行一次查询， ==
echo == 如发现其停止，则立即启动。 ==
echo ========================================
echo.
echo 此脚本监测的服务是：%_task%,端口:%_port%
echo.


:checkstart
::通过进程的端口是否正在被监听检测
netstat -ano | find "0.0.0.0:%_port%" >nul 2>nul && goto checkag|| goto startsvr

:startsvr
echo %time%
echo ********程序开始启动********
echo 程序重新启动于 %time% ,请检查系统日志 >> restart_service.log
echo cd /d %cd% > %_des%
echo start %_task% >> %_des%
echo exit >> %_des%
start %_des%
set/p=.<nul
for /L %%i in (1 1 10) do set /p a=.<nul&ping.exe /n 2 127.0.0.1>nul
echo .
echo Wscript.Sleep WScript.Arguments(0) >%tmp%\delay.vbs
cscript //b //nologo %tmp%\delay.vbs 10000
del %_des% /Q

echo ********程序启动完成********
goto checkstart

:checkag
echo %time% 程序运行正常,10秒后继续检查.. 
echo Wscript.Sleep WScript.Arguments(0) >%tmp%\delay.vbs
cscript //b //nologo %tmp%\delay.vbs 10000

goto checkstart