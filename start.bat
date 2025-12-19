@echo off
REM TanGo 项目启动脚本 (Windows版本)
REM 功能：读取 .env 配置文件，同时启动前端和后端服务

setlocal enabledelayedexpansion

REM 项目根目录
set "ROOT_DIR=%~dp0"
set "BACKEND_DIR=%ROOT_DIR%backend"
set "FRONTEND_DIR=%ROOT_DIR%frontend"
set "ENV_FILE=%ROOT_DIR%.env"

echo ========================================
echo   TanGo 项目启动脚本
echo ========================================
echo.

REM 检查 .env 文件
if not exist "%ENV_FILE%" (
    echo 警告: .env 文件不存在
    if exist "%ROOT_DIR%.env.example" (
        echo 提示: 发现 .env.example 文件，是否复制为 .env? (Y/N)
        set /p answer=
        if /i "!answer!"=="Y" (
            copy "%ROOT_DIR%.env.example" "%ENV_FILE%"
            echo 已复制 .env.example 为 .env
            echo 请编辑 .env 文件，填入实际的配置值
            pause
            exit /b 1
        )
    )
    echo 错误: 需要 .env 配置文件
    pause
    exit /b 1
)
echo [OK] 找到 .env 配置文件

REM 加载 .env 文件
if exist "%ENV_FILE%" (
    for /f "usebackq tokens=1,* delims==" %%a in ("%ENV_FILE%") do (
        set "line=%%a"
        if not "!line:~0,1!"=="#" (
            if not "!line!"=="" (
                set "%%a"
            )
        )
    )
    echo [OK] 已加载 .env 配置
)

REM 检查依赖
echo 检查依赖...
where go >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 Go，请先安装 Go 1.21+
    pause
    exit /b 1
)
echo [OK] Go 已安装

where node >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 Node.js，请先安装 Node.js
    pause
    exit /b 1
)
echo [OK] Node.js 已安装

where npm >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 npm，请先安装 npm
    pause
    exit /b 1
)
echo [OK] npm 已安装

REM 安装前端依赖（如果需要）
if not exist "%FRONTEND_DIR%\node_modules" (
    echo 前端依赖未安装，正在安装...
    cd /d "%FRONTEND_DIR%"
    call npm install
    cd /d "%ROOT_DIR%"
    echo [OK] 前端依赖安装完成
)

REM 启动后端服务
echo.
echo 启动后端服务...
cd /d "%BACKEND_DIR%"

if not defined BACKEND_PORT set "BACKEND_PORT=8877"
if not defined BACKEND_HOST set "BACKEND_HOST=0.0.0.0"

echo 后端服务地址: %BACKEND_HOST%:%BACKEND_PORT%

start "TanGo Backend" cmd /c "go run explore.go -f etc/explore.yaml > %ROOT_DIR%\backend.log 2>&1"
timeout /t 3 /nobreak >nul
echo [OK] 后端服务已启动

cd /d "%ROOT_DIR%"

REM 启动前端服务
echo.
echo 启动前端服务...
cd /d "%FRONTEND_DIR%"

if not defined FRONTEND_PORT set "FRONTEND_PORT=3000"

echo 前端服务地址: http://localhost:%FRONTEND_PORT%

start "TanGo Frontend" cmd /c "npm run dev > %ROOT_DIR%\frontend.log 2>&1"
timeout /t 5 /nobreak >nul
echo [OK] 前端服务已启动

cd /d "%ROOT_DIR%"

REM 显示服务信息
echo.
echo ========================================
echo   TanGo 服务已启动
echo ========================================
echo 后端服务: http://%BACKEND_HOST%:%BACKEND_PORT%
echo 前端服务: http://localhost:%FRONTEND_PORT%
echo.
echo 查看日志:
echo   后端: type %ROOT_DIR%\backend.log
echo   前端: type %ROOT_DIR%\frontend.log
echo.
echo 按任意键关闭服务窗口...
echo ========================================
echo.

pause
