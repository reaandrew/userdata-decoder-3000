param (
    [string]$Executable,
    [string]$Os
)

# Execute the command with parameters
& $Executable --provider aws -v --output-dir ./output/$Os/
$files = Get-ChildItem -Path .\output\$Os -Recurse -File | Select-Object -ExpandProperty FullName
$fileCount = $files.Count
Write-Output "Number of files: $fileCount"

# Read required files from a list file
$requiredFiles = Get-Content -Path .\scripts\ci\expected_files.txt

# Print out required files for debugging
Write-Output "Required files are:"
$requiredFiles | ForEach-Object { Write-Output $_ }

# Check for each required file and exit with an error if any are missing
$foundFiles = $files | Where-Object { $requiredFiles -contains $_.Name }

if ($foundFiles.Count -ne $requiredFiles.Count) {
    Write-Error "Error: Required files are missing!"
    exit 1
}

if ($fileCount -ne 6) {
    Write-Error "Error: Incorrect number of files found! Expected 6, found $fileCount"
    exit 1
}