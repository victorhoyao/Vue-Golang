# Get fresh token
$loginBody = @{
    userName = "testadmin"
    passWord = "test123"
} | ConvertTo-Json

$headers = @{
    "Content-Type" = "application/json"
}

Write-Host "Getting fresh token..."

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8036/login" -Method POST -Body $loginBody -Headers $headers -UseBasicParsing
    $data = $response.Content | ConvertFrom-Json
    
    if ($data.code -eq 200) {
        Write-Host "Login successful!"
        Write-Host "Token: $($data.data.token)"
        
        # Now test TaskItem creation with fresh token
        $taskHeaders = @{
            "Content-Type" = "application/json"
            "token" = $data.data.token
        }
        
        $taskBody = @{
            taskId = "TEST002"
            taskName = "Test Task 2"
            supplier = "Test Supplier 2"
            userTypes = "premium"
            priceReal = 10.50
            priceProtocol = 12.00
            priceManual = 15.00
            selectionMode = "Single-select"
            status = "active"
        } | ConvertTo-Json
        
        Write-Host "Testing TaskItem creation with fresh token..."
        $taskResponse = Invoke-WebRequest -Uri "http://localhost:8036/TaskItem/add" -Method POST -Body $taskBody -Headers $taskHeaders -UseBasicParsing
        Write-Host "TaskItem Response Status: $($taskResponse.StatusCode)"
        Write-Host "TaskItem Response Content: $($taskResponse.Content)"
        
    } else {
        Write-Host "Login failed: $($data.msg)"
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)"
} 