# URP Specification Compliance Test Suite
param([switch]$Verbose)

$ErrorActionPreference = "Stop"
$script:totalTests = 0
$script:passedTests = 0
$script:failedTests = 0

function Write-TestHeader {
    param($text)
    Write-Host "`n===============================================" -ForegroundColor Cyan
    Write-Host $text -ForegroundColor Cyan
    Write-Host "===============================================" -ForegroundColor Cyan
}

function Write-TestResult {
    param($passed, $message)
    $script:totalTests++
    if ($passed) {
        $script:passedTests++
        Write-Host "  PASS: $message" -ForegroundColor Green
    } else {
        $script:failedTests++
        Write-Host "  FAIL: $message" -ForegroundColor Red
    }
}

function Test-LogFormat {
    param([string]$LogFile, [string]$Type)
    
    Write-TestHeader "Testing $Type Log Format Compliance"
    
    if (!(Test-Path $LogFile)) {
        Write-TestResult $false "$Type log file not found"
        return
    }
    
    $lines = Get-Content $LogFile
    $logLines = $lines | Where-Object { $_ -match '^\w+\s+\w+\s+[\d.]+\s+\w+\s+\d+\s+\d+' }
    
    # Test log format
    $formatCorrect = $true
    foreach ($line in $logLines) {
        if ($line -notmatch '^(snd|rcv)\s+(ok|drp|cor)\s+[\d.]+\s+(DATA|ACK|SYN|FIN)\s+\d+\s+\d+$') {
            Write-TestResult $false "Invalid log format: $line"
            $formatCorrect = $false
            break
        }
    }
    if ($formatCorrect -and $logLines.Count -gt 0) {
        Write-TestResult $true "Log format matches spec (checked $($logLines.Count) lines)"
    }
}

function Test-Statistics {
    param([string]$LogFile, [string]$Type)
    
    Write-TestHeader "Testing $Type Statistics Compliance"
    
    if (!(Test-Path $LogFile)) {
        Write-TestResult $false "$Type log file not found"
        return
    }
    
    $content = Get-Content $LogFile -Raw
    
    if ($Type -eq "Sender") {
        $requiredStats = @(
            "Original data sent:",
            "Total data sent:",
            "Original segments sent:",
            "Total segments sent:",
            "Timeout retransmissions:",
            "Fast retransmissions:",
            "Duplicate acks received:",
            "Corrupted acks discarded:",
            "PLC forward segments dropped:",
            "PLC forward segments corrupted:",
            "PLC reverse segments dropped:",
            "PLC reverse segments corrupted:"
        )
    } else {
        $requiredStats = @(
            "Original data received:",
            "Total data received:",
            "Original segments received:",
            "Total segments received:",
            "Corrupted segments discarded:",
            "Duplicate segments received:",
            "Total acks sent:",
            "Duplicate acks sent:"
        )
    }
    
    foreach ($stat in $requiredStats) {
        $found = $content -match [regex]::Escape($stat)
        Write-TestResult $found "Statistic present: $stat"
    }
}

# Main execution
Write-Host "URP Specification Compliance Test Suite" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow

# Build
Write-Host "`nBuilding project..." -ForegroundColor Cyan
go build -o sender.exe .\cmd\sender
go build -o receiver.exe .\cmd\receiver

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

# Clean up
Remove-Item "data\sender_log.txt" -Force -ErrorAction SilentlyContinue
Remove-Item "data\receiver_log.txt" -Force -ErrorAction SilentlyContinue
Remove-Item "data\received_file.txt" -Force -ErrorAction SilentlyContinue

# Create test file
if (!(Test-Path "data")) {
    New-Item -ItemType Directory -Path "data" | Out-Null
}

$testContent = "This is a test file for URP compliance testing. " * 70
Set-Content -Path "data\test_compliance.txt" -Value $testContent -NoNewline

# Run test
Write-Host "`nRunning test transfer..." -ForegroundColor Cyan

$receiverArgs = "-receiver_port=5000 -sender_port=6000 -file=data/received_file.txt -max_win=3000"
$receiver = Start-Process -FilePath ".\receiver.exe" -ArgumentList $receiverArgs.Split() -NoNewWindow -PassThru
Start-Sleep -Seconds 2

$senderArgs = "-sender_port=6000 -receiver_host=127.0.0.1 -receiver_port=5000 -file=data/test_compliance.txt -max_win=3000 -rto=100 -flp=0.05 -rlp=0.05 -fcp=0.05 -rcp=0.05"
$senderProc = Start-Process -FilePath ".\sender.exe" -ArgumentList $senderArgs.Split() -NoNewWindow -PassThru -Wait

Start-Sleep -Seconds 3
if (!$receiver.HasExited) {
    Stop-Process -Id $receiver.Id -Force -ErrorAction SilentlyContinue
}
Start-Sleep -Milliseconds 500

# Run tests
Test-LogFormat "data\sender_log.txt" "Sender"
Test-LogFormat "data\receiver_log.txt" "Receiver"
Test-Statistics "data\sender_log.txt" "Sender"
Test-Statistics "data\receiver_log.txt" "Receiver"

# Summary
Write-TestHeader "TEST SUMMARY"
Write-Host "`nTotal Tests: $script:totalTests" -ForegroundColor Cyan
Write-Host "Passed: $script:passedTests" -ForegroundColor Green
Write-Host "Failed: $script:failedTests" -ForegroundColor Red

if ($script:failedTests -gt 0) {
    exit 1
} else {
    exit 0
}

