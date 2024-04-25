param (
    [string]$Executable,
    [string]$Os
)

# Execute the command with parameters
& $Executable --provider aws -v --output-dir ./output/$Os/
$files = Get-ChildItem -Path .\output\$Os -Recurse -File | Select-Object -ExpandProperty Name
Write-Output "Files found: $($files -join ', ')"

# Read required files from a list file
$requiredFiles = Get-Content -Path .\scripts\ci\expected_files.txt

# Print out required files for debugging
Write-Output "Required files are:"
$requiredFiles | ForEach-Object { Write-Output $_ }

# Check for each required file and exit with an error if any are missing
foreach ($file in $requiredFiles) {
    if (-not ($files -contains $file)) {
        Write-Error "Error: Can't find $file"
        exit 1
    }
}
