param (
    [string]$Executable,
    [string]$Os
)

# Execute the command with parameters
& $Executable --provider aws -v --output-dir ./output/$Os/
$files = Get-ChildItem -Path .\output\$Os -Recurse -File | Select-Object -ExpandProperty FullName
$fileCount = $files.Count
Write-Output "Number of files: $fileCount"

# Read required files from a list file and handle potential CR characters
$requiredFiles = Get-Content -Path .\scripts\ci\expected_files.txt | ForEach-Object { $_.TrimEnd("`r") }

# Normalize file names to a standard format if needed, e.g., lowercasing or similar
$foundFiles = $files | Where-Object { $requiredFiles -contains ($_.Name.ToLower()) }


if ($foundFiles.Count -ne $requiredFiles.Count) {
    Write-Error "Error: Required files are missing!"
    exit 1
}

if ($fileCount -ne 6) {
    Write-Error "Error: Incorrect number of files found! Expected 6, found $fileCount"
    exit 1
}