@echo off
REM TanGo 项目停止脚本 (Windows版本)
REM 功能：停止所有运行中的前后端服务

setlocal enabledelayedexpansion

REM 项目根目录
set "ROOT_DIR=%~dp0"
set "BACKEND_DIR=%ROOT_DIR%backend"
set "FRONTEND_DIR=%ROOT_DIR%frontend"

REM 加载 .env 文件获取端口配置
set "ENV_FILE=%ROOT_DIR%.env"
if exist "%ENV_FILE%" (
    for /f "usebackq tokens=1,* delims==" %%a in ("%ENV_FILE%") do (
        set "line=%%a"
        if not "!line:~0,1!"=="#" (
            if not "!line!"=="" (
                set "%%a"
            )
        )
    )
)

if not defined BACKEND_PORT set "BACKEND_PORT=8877"
if not defined FRONTEND_PORT set "FRONTEND_PORT=3000"

echo ========================================
echo   TanGo 服务停止脚本
echo ========================================
echo.

REM 停止后端服务
echo 停止后端服务...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%BACKEND_PORT%"') do (
    set "PID=%%a"
    if not "!PID!"=="" (
        echo 正在停止后端进程 (PID: !PID!)...
        taskkill /F /PID !PID! >nul 2>&1
    )
)

REM 停止前端服务
echo 停止前端服务...
for /f "tokens=5" %%a in ('netstat -ano ^| findstr ":%FRONTEND_PORT%"') do (
    set "PID=%%a"
    if not "!PID!"=="" (
        echo 正在停止前端进程 (PID: !PID!)...
        taskkill /F /PID !PID! >nul 2>&1
    )
)

REM 停止所有相关进程
taskkill /F /FI "WINDOWTITLE eq TanGo Backend*" >nul 2>&1
taskkill /F /FI "WINDOWTITLE eq TanGo Frontend*" >nul 2>&1

echo.
echo ========================================
echo   所有服务已停止
echo ========================================
echo.
pause
