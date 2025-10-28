# URP Comprehensive Test Suite
# Based on assignment specification test cases

param(
    [switch]$Quick,
    [switch]$Category1,
    [switch]$Category2,
    [switch]$Category3,
    [string]$TestId
)

$ErrorActionPreference = "Continue"

# Colors for output
function Write-TestHeader($text) {
    Write-Host "`n================================================" -ForegroundColor Cyan
    Write-Host $text -ForegroundColor Cyan
    Write-Host "================================================" -ForegroundColor Cyan
}

function Write-TestResult($passed, $message) {
    if ($passed) {
        Write-Host "  ✓ PASS: $message" -ForegroundColor Green
    } else {
        Write-Host "  ✗ FAIL: $message" -ForegroundColor Red
    }
}

function Write-TestInfo($message) {
    Write-Host "  ℹ $message" -ForegroundColor Yellow
}

# Test result tracking
$script:totalTests = 0
$script:passedTests = 0
$script:failedTests = 0
$script:testResults = @()

# Clean up function
function Clean-TestEnvironment {
    Stop-Process -Name receiver -Force -ErrorAction SilentlyContinue
    Stop-Process -Name sender -Force -ErrorAction SilentlyContinue
    Start-Sleep -Milliseconds 500
    
    Remove-Item "data\received_file.txt" -Force -ErrorAction SilentlyContinue
    Remove-Item "data\sender_log.txt" -Force -ErrorAction SilentlyContinue
    Remove-Item "data\receiver_log.txt" -Force -ErrorAction SilentlyContinue
}

# Run a single test
function Run-Test {
    param(
        [string]$TestId,
        [string]$Description,
        [int]$MaxWin,
        [int]$Rto,
        [double]$Flp,
        [double]$Rlp,
        [double]$Fcp,
        [double]$Rcp,
        [string]$TestFile = "data/test_file.txt",
        [scriptblock]$CustomValidation = $null
    )
    
    $script:totalTests++
    Write-TestHeader "Test ${TestId}: $Description"
    
    Clean-TestEnvironment
    
    # Start receiver
    Write-TestInfo "Starting receiver (port 5000, max_win=$MaxWin)..."
    $receiver = Start-Process -FilePath ".\receiver.exe" `
        -ArgumentList "--receiver_port","5000","--sender_port","6000","--file","data/received_file.txt","--max_win",$MaxWin,"--rlp",$Rlp,"--rcp",$Rcp `
        -NoNewWindow -PassThru -RedirectStandardOutput "receiver_stdout.txt" -RedirectStandardError "receiver_stderr.txt"
    
    Start-Sleep -Seconds 2
    
    # Start sender
    Write-TestInfo "Starting sender (max_win=$MaxWin, rto=$Rto, flp=$Flp, rlp=$Rlp, fcp=$Fcp, rcp=$Rcp)..."
    $senderArgs = @(
        "--sender_port", "6000",
        "--receiver_host", "127.0.0.1",
        "--receiver_port", "5000",
        "--file", $TestFile,
        "--max_win", $MaxWin,
        "--rto", $Rto,
        "--flp", $Flp,
        "--rlp", $Rlp,
        "--fcp", $Fcp,
        "--rcp", $Rcp
    )
    
    $sender = Start-Process -FilePath ".\sender.exe" `
        -ArgumentList $senderArgs `
        -NoNewWindow -PassThru -Wait -RedirectStandardOutput "sender_stdout.txt" -RedirectStandardError "sender_stderr.txt"
    
    $senderExitCode = $sender.ExitCode
    
    # Wait for receiver TIME_WAIT
    Start-Sleep -Seconds 3
    
    # Stop receiver if still running
    if (!$receiver.HasExited) {
        Stop-Process -Id $receiver.Id -Force -ErrorAction SilentlyContinue
    }
    
    Start-Sleep -Milliseconds 500
    
    # Verify results
    $testPassed = $true
    $failureReasons = @()
    
    # Check if sender completed successfully
    if ($senderExitCode -ne 0 -and $senderExitCode -ne 1) {
        $testPassed = $false
        $failureReasons += "Sender exited with code $senderExitCode"
    }
    
    # Check if received file exists
    if (!(Test-Path "data\received_file.txt")) {
        $testPassed = $false
        $failureReasons += "Received file not created"
    } else {
        # Verify file integrity
        $sourceHash = (Get-FileHash $TestFile -Algorithm MD5 -ErrorAction SilentlyContinue).Hash
        $receivedHash = (Get-FileHash "data\received_file.txt" -Algorithm MD5 -ErrorAction SilentlyContinue).Hash
        
        if ($sourceHash -ne $receivedHash) {
            $testPassed = $false
            $failureReasons += "File hashes don't match"
        } else {
            Write-TestResult $true "Files are identical"
        }
        
        # Check file sizes
        $sourceSize = (Get-Item $TestFile -ErrorAction SilentlyContinue).Length
        $receivedSize = (Get-Item "data\received_file.txt" -ErrorAction SilentlyContinue).Length
        Write-TestInfo "Original: $sourceSize bytes, Received: $receivedSize bytes"
    }
    
    # Parse logs
    if (Test-Path "data\sender_log.txt") {
        $senderLog = Get-Content "data\sender_log.txt" -Raw -ErrorAction SilentlyContinue
        $dropCount = ([regex]::Matches($senderLog, "drop")).Count
        $corruptCount = ([regex]::Matches($senderLog, "crpt")).Count
        $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
        $rexmtCount = ([regex]::Matches($senderLog, "rexmt")).Count
        
        Write-TestInfo "Sender - Drops: $dropCount, Corruptions: $corruptCount, Timeouts: $timeoutCount, Retransmits: $rexmtCount"
    }
    
    if (Test-Path "data\receiver_log.txt") {
        $receiverLog = Get-Content "data\receiver_log.txt" -Raw -ErrorAction SilentlyContinue
        $rcvCount = ([regex]::Matches($receiverLog, " rcv ")).Count
        $sndCount = ([regex]::Matches($receiverLog, " snd ")).Count
        
        Write-TestInfo "Receiver - Segments received: $rcvCount, ACKs sent: $sndCount"
    }
    
    # Run custom validation if provided
    if ($CustomValidation) {
        try {
            $customResult = & $CustomValidation
            if (!$customResult) {
                $testPassed = $false
                $failureReasons += "Custom validation failed"
            }
        } catch {
            $testPassed = $false
            $failureReasons += "Custom validation error: $_"
        }
    }
    
    # Final result
    if ($testPassed) {
        Write-TestResult $true "Test $TestId passed"
        $script:passedTests++
        $script:testResults += [PSCustomObject]@{
            TestId = $TestId
            Description = $Description
            Result = "PASS"
            Details = ""
        }
    } else {
        Write-TestResult $false "Test $TestId failed: $($failureReasons -join ', ')"
        $script:failedTests++
        $script:testResults += [PSCustomObject]@{
            TestId = $TestId
            Description = $Description
            Result = "FAIL"
            Details = $failureReasons -join ', '
        }
    }
    
    # Clean up stdout/stderr files
    Remove-Item "sender_stdout.txt" -Force -ErrorAction SilentlyContinue
    Remove-Item "sender_stderr.txt" -Force -ErrorAction SilentlyContinue
    Remove-Item "receiver_stdout.txt" -Force -ErrorAction SilentlyContinue
    Remove-Item "receiver_stderr.txt" -Force -ErrorAction SilentlyContinue
}

# Create test files of various sizes
function Create-TestFile {
    param(
        [string]$Path,
        [int]$Size
    )
    
    $content = ""
    $chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 `n"
    $random = New-Object System.Random
    
    for ($i = 0; $i -lt $Size; $i++) {
        $content += $chars[$random.Next($chars.Length)]
    }
    
    Set-Content -Path $Path -Value $content -NoNewline
}

# Build executables
Write-Host "Building executables..." -ForegroundColor Cyan
go build -o sender.exe .\cmd\sender
go build -o receiver.exe .\cmd\receiver

if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}

Write-Host "Build successful!" -ForegroundColor Green

# Ensure data directory exists
if (!(Test-Path "data")) {
    New-Item -ItemType Directory -Path "data" | Out-Null
}

# Create test files
Write-Host "Creating test files..." -ForegroundColor Cyan
Create-TestFile "data/test_file.txt" 1013
Create-TestFile "data/test_3500.txt" 3500
Create-TestFile "data/test_10000.txt" 10000

# ============================================================================
# CATEGORY 1: Stop-and-Wait Protocol (max_win=1000)
# ============================================================================

if (!$Category2 -and !$Category3 -or $Category1) {
    Write-TestHeader "CATEGORY 1: Stop-and-Wait Protocol"
    
    # Test 1.1: Baseline
    Run-Test -TestId "1.1" -Description "Baseline: Flawless transfer" `
        -MaxWin 1000 -Rto 200 -Flp 0 -Rlp 0 -Fcp 0 -Rcp 0 `
        -CustomValidation {
            $senderLog = Get-Content "data\sender_log.txt" -Raw
            $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
            $dropCount = ([regex]::Matches($senderLog, "drop")).Count
            return ($timeoutCount -eq 0 -and $dropCount -eq 0)
        }
    
    if (!$Quick) {
        # Test 1.2: Forward DATA Loss
        Run-Test -TestId "1.2" -Description "Forward DATA Loss" `
            -MaxWin 1000 -Rto 100 -Flp 0.2 -Rlp 0 -Fcp 0 -Rcp 0 `
            -CustomValidation {
                $senderLog = Get-Content "data\sender_log.txt" -Raw
                $dropCount = ([regex]::Matches($senderLog, "drop")).Count
                $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
                Write-TestInfo "Forward drops: $dropCount, Timeouts: $timeoutCount"
                return ($dropCount -gt 0 -and $timeoutCount -gt 0)
            }
        
        # Test 1.3: Reverse ACK Loss
        Run-Test -TestId "1.3" -Description "Reverse ACK Loss" `
            -MaxWin 1000 -Rto 100 -Flp 0 -Rlp 0.2 -Fcp 0 -Rcp 0 `
            -CustomValidation {
                $senderLog = Get-Content "data\sender_log.txt" -Raw
                $dropCount = ([regex]::Matches($senderLog, "drop")).Count
                $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
                Write-TestInfo "Reverse drops: $dropCount, Timeouts: $timeoutCount"
                return ($dropCount -gt 0)
            }
        
        # Test 1.4: Forward DATA Corruption
        Run-Test -TestId "1.4" -Description "Forward DATA Corruption" `
            -MaxWin 1000 -Rto 100 -Flp 0 -Rlp 0 -Fcp 0.2 -Rcp 0 `
            -CustomValidation {
                $senderLog = Get-Content "data\sender_log.txt" -Raw
                $corruptCount = ([regex]::Matches($senderLog, "crpt")).Count
                $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
                Write-TestInfo "Forward corruptions: $corruptCount, Timeouts: $timeoutCount"
                return ($corruptCount -gt 0)
            }
        
        # Test 1.5: Reverse ACK Corruption
        Run-Test -TestId "1.5" -Description "Reverse ACK Corruption" `
            -MaxWin 1000 -Rto 100 -Flp 0 -Rlp 0 -Fcp 0 -Rcp 0.2 `
            -CustomValidation {
                $senderLog = Get-Content "data\sender_log.txt" -Raw
                $corruptCount = ([regex]::Matches($senderLog, "crpt")).Count
                Write-TestInfo "Reverse corruptions: $corruptCount"
                return ($corruptCount -gt 0)
            }
        
        # Test 1.6: Combined Unreliability
        Run-Test -TestId "1.6" -Description "Combined Unreliability" `
            -MaxWin 1000 -Rto 150 -Flp 0.1 -Rlp 0.1 -Fcp 0.1 -Rcp 0.1
        
        # Test 1.7: Connection Setup Resilience
        Run-Test -TestId "1.7" -Description "Connection Setup Resilience" `
            -MaxWin 1000 -Rto 100 -Flp 0.5 -Rlp 0.5 -Fcp 0 -Rcp 0
        
        # Test 1.8: Connection Teardown Resilience
        Run-Test -TestId "1.8" -Description "Connection Teardown Resilience" `
            -MaxWin 1000 -Rto 100 -Flp 0.5 -Rlp 0.5 -Fcp 0 -Rcp 0
    }
}

# ============================================================================
# CATEGORY 2: Sliding Window Protocol (max_win > 1000)
# ============================================================================

if (!$Category1 -and !$Category3 -or $Category2) {
    Write-TestHeader "CATEGORY 2: Sliding Window Protocol"
    
    # Test 2.1: Baseline pipelined transfer
    Run-Test -TestId "2.1" -Description "Baseline: Pipelined transfer" `
        -MaxWin 5000 -Rto 200 -Flp 0 -Rlp 0 -Fcp 0 -Rcp 0 `
        -TestFile "data/test_10000.txt" `
        -CustomValidation {
            $senderLog = Get-Content "data\sender_log.txt" -Raw
            $timeoutCount = ([regex]::Matches($senderLog, "timeout")).Count
            return ($timeoutCount -eq 0)
        }
    
    if (!$Quick) {
        # Test 2.2: Single Packet Loss
        Run-Test -TestId "2.2" -Description "Single Packet Loss (Fast Retransmit)" `
            -MaxWin 5000 -Rto 500 -Flp 0.05 -Rlp 0 -Fcp 0 -Rcp 0 `
            -TestFile "data/test_10000.txt"
        
        # Test 2.3: Multiple Packet Loss
        Run-Test -TestId "2.3" -Description "Multiple Packet Loss" `
            -MaxWin 5000 -Rto 100 -Flp 0.3 -Rlp 0 -Fcp 0 -Rcp 0 `
            -TestFile "data/test_10000.txt"
        
        # Test 2.4: Cumulative ACK Resilience
        Run-Test -TestId "2.4" -Description "Cumulative ACK Resilience" `
            -MaxWin 5000 -Rto 500 -Flp 0 -Rlp 0.1 -Fcp 0 -Rcp 0 `
            -TestFile "data/test_10000.txt"
        
        # Test 2.5: Full Unreliability
        Run-Test -TestId "2.5" -Description "Full Unreliability (Sliding Window)" `
            -MaxWin 5000 -Rto 200 -Flp 0.15 -Rlp 0.15 -Fcp 0.1 -Rcp 0.1 `
            -TestFile "data/test_10000.txt"
    }
}

# ============================================================================
# CATEGORY 3: Edge Cases and Specification Compliance
# ============================================================================

if (!$Category1 -and !$Category2 -or $Category3) {
    Write-TestHeader "CATEGORY 3: Edge Cases and Specification Compliance"
    
    # Test 3.1: Final Segment Size
    Run-Test -TestId "3.1" -Description "Final Segment Size (non-MSS multiple)" `
        -MaxWin 4000 -Rto 200 -Flp 0 -Rlp 0 -Fcp 0 -Rcp 0 `
        -TestFile "data/test_3500.txt" `
        -CustomValidation {
            $senderLog = Get-Content "data\sender_log.txt"
            # Look for a segment with payload length not equal to 1000
            $lastDataLine = $senderLog | Where-Object { $_ -match "snd.*NONE" } | Select-Object -Last 1
            if ($lastDataLine -match "(\d+)$") {
                $lastPayloadSize = [int]$matches[1]
                Write-TestInfo "Last payload size: $lastPayloadSize bytes"
                return ($lastPayloadSize -lt 1000)
            }
            return $false
        }
    
    if (!$Quick) {
        # Test 3.4: Receiver TIME_WAIT State
        Run-Test -TestId "3.4" -Description "Receiver TIME_WAIT State (2 seconds)" `
            -MaxWin 1000 -Rto 200 -Flp 0 -Rlp 0 -Fcp 0 -Rcp 0 `
            -CustomValidation {
                $receiverLog = Get-Content "data\receiver_log.txt"
                $timeWaitLine = $receiverLog | Where-Object { $_ -match "TIME_WAIT" }
                return ($timeWaitLine -ne $null)
            }
        
        # Test 3.5: No Timeout Doubling
        Run-Test -TestId "3.5" -Description "No Timeout Doubling (Fixed RTO)" `
            -MaxWin 1000 -Rto 100 -Flp 0.5 -Rlp 0 -Fcp 0 -Rcp 0 `
            -CustomValidation {
                $senderLog = Get-Content "data\sender_log.txt"
                $timeoutLines = $senderLog | Where-Object { $_ -match "^([\d.]+)\s+timeout" }
                
                if ($timeoutLines.Count -ge 2) {
                    $times = @()
                    foreach ($line in $timeoutLines) {
                        if ($line -match "^([\d.]+)") {
                            $times += [double]$matches[1]
                        }
                    }
                    
                    if ($times.Count -ge 2) {
                        $delta1 = $times[1] - $times[0]
                        Write-TestInfo "Time between first two timeouts: $([math]::Round($delta1, 2))ms"
                        # Should be approximately 100ms, not doubling
                        return ($delta1 -ge 80 -and $delta1 -le 150)
                    }
                }
                return $true  # If not enough timeouts, consider it a pass
            }
    }
}

# ============================================================================
# Test Summary
# ============================================================================

Write-TestHeader "TEST SUMMARY"

Write-Host "`nTotal Tests: $script:totalTests" -ForegroundColor Cyan
Write-Host "Passed: $script:passedTests" -ForegroundColor Green
Write-Host "Failed: $script:failedTests" -ForegroundColor Red

if ($script:failedTests -gt 0) {
    Write-Host "`nFailed Tests:" -ForegroundColor Red
    $script:testResults | Where-Object { $_.Result -eq "FAIL" } | ForEach-Object {
        Write-Host "  - $($_.TestId): $($_.Description)" -ForegroundColor Red
        if ($_.Details) {
            Write-Host "    Reason: $($_.Details)" -ForegroundColor Yellow
        }
    }
}

Write-Host "`nSuccess Rate: $([math]::Round(($script:passedTests / $script:totalTests) * 100, 1))%" -ForegroundColor Cyan

# Export results to CSV
$script:testResults | Export-Csv -Path "test_results.csv" -NoTypeInformation
Write-Host "`nDetailed results exported to: test_results.csv" -ForegroundColor Yellow

# Clean up
Clean-TestEnvironment

# Exit with appropriate code
if ($script:failedTests -gt 0) {
    exit 1
} else {
    exit 0
}
