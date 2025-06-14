# Test TaskItem API
$token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjgsImlhdCI6MTc0OTg1NTc0NywiaXNzIjoiYXV0aG9yIiwic3ViIjoidXNlciB0b2tlbiJ9.LQGJ40Dsjd1ijNlwAnWccCRoAXq4"

$headers = @{
    "Content-Type" = "application/json"
    "token" = $token
}

$body = @{
    taskId = "TEST001"
    taskName = "Test Task"
    supplier = "Test Supplier"
    userTypes = "premium"
    priceReal = 10.50
    priceProtocol = 12.00
    priceManual = 15.00
    selectionMode = "Single-select"
    status = "active"
} | ConvertTo-Json

Write-Host "Testing TaskItem creation..."
Write-Host "Body: $body"

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8036/TaskItem/add" -Method POST -Body $body -Headers $headers -UseBasicParsing
    Write-Host "Response Status: $($response.StatusCode)"
    Write-Host "Response Content: $($response.Content)"
} catch {
    Write-Host "Error: $($_.Exception.Message)"
    Write-Host "Response: $($_.Exception.Response)"
} 