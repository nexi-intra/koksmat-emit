param (
  $requiredVars = @(
    'DEVICEKPIBLOB',
    'GRAPH_APPID',
    'GRAPH_APPSECRET',
    'GRAPH_APPDOMAIN',
    'GRAPH_ACCESSTOKEN',
    'POSTGRES_SERVER'
  
  )
)

# Enumerate all environment variables
$allEnv = Get-ChildItem Env:

# Store environment variables in a hashtable for easy lookup
$envHash = @{}
$allEnv | ForEach-Object {
  $envHash[$_.Name] = $_.Value
}
  
# Check which required variables are missing or empty
$missingVars = $requiredVars | Where-Object { -not $envHash[$_] -or $envHash[$_] -eq "" }
  
if ($missingVars.Count -gt 0) {
  Write-Host "The following required environment variables are missing or empty:"
  $missingVars | ForEach-Object { Write-Host "- $_" }
  throw "Missing or empty environment variables."
  exit 1
}


# Check for variables with leading or trailing whitespace
$whitespaceVars = $requiredVars | Where-Object { $envHash[$_] -match '^\s|\s$' }

if ($whitespaceVars.Count -gt 0) {
  Write-Host "The following environment variables have leading or trailing whitespace:"
  $whitespaceVars | ForEach-Object { Write-Host "- $_" }
  throw "Environment variables with leading or trailing whitespace found."
  exit
}

Write-Host "All required environment variables are valid."