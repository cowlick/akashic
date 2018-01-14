$commands = @(
    "akashic",
    "akashic-new"
)

if (Test-Path .\release) {
    rm .\release -Recurse -Force
}

$commands | foreach {
    gox -verbose -output "./release/${_}_{{.OS}}_{{.Arch}}" ./${_}
}

cd .\release

ls | where {$_.Name -match '^akashic-new*'} | foreach {
    $bin = $_.Name -replace "-new",""
    $list  = New-Object 'System.Collections.Generic.List[System.String]'
    $commands | foreach {
        if ($bin -match "windows") {
            $command = "${_}.exe"
        } else {
            $command = $_
        }
        $list.Add($command)
        mv ($bin -replace "akashic", $_) $command
    }
    $list | Compress-Archive -DestinationPath "${bin}.zip"
    $list | rm
}
