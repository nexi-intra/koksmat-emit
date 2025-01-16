
# Define the path to the JSON file
$jsonFilePath = join-path "$PSScriptRoot" ".." "./.koksmat" "koksmat.json"

# Check if the JSON file exists
if (-Not (Test-Path -Path $jsonFilePath)) {
  Write-Error "The file '$jsonFilePath' does not exist."
  exit 1
}

try {
  # Read the content of the JSON file
  $jsonContent = Get-Content -Path $jsonFilePath -Raw | ConvertFrom-Json

  # Ensure the 'version' property exists
  if (-Not $jsonContent.PSObject.Properties.Match("version")) {
    Write-Error "The JSON does not contain a 'version' property."
    exit 1
  }

  # Increment the 'build' number by 1
  if ($jsonContent.version.PSObject.Properties.Match("build")) {
    $jsonContent.version.build += 1
  }
  else {
    Write-Error "The 'version' object does not contain a 'build' property."
    exit 1
  }

  # Increment the 'patch' number by 1
  if ($jsonContent.version.PSObject.Properties.Match("patch")) {
    $jsonContent.version.patch += 1
  }
  else {
    Write-Error "The 'version' object does not contain a 'patch' property."
    exit 1
  }

  # Convert the updated object back to JSON with indentation for readability
  $updatedJson = $jsonContent | ConvertTo-Json -Depth 4 -Compress:$false

  # Write the updated JSON back to the file
  Set-Content -Path $jsonFilePath -Value $updatedJson

  Write-Host "Successfully incremented 'build' and 'patch' numbers in '$jsonFilePath'."
}
catch {
  Write-Error "An error occurred: $_"
  exit 1
}

$tag = "v$($jsonContent.version.major).$($jsonContent.version.minor).$($jsonContent.version.patch).$($jsonContent.version.build)"

Write-Host "Releasing $tag"

gh release create $tag  --generate-notes  

