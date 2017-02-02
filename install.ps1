$pwd = Get-Location
$servicePath = "${pwd}\haproxy-dynagent.exe"
Write-Host "Registering service for ${servicePath}"
New-Service -Name "HAProxyDynAgent" -DisplayName "HAProxy DynAgent" -BinaryPathName $servicePath -StartupType Automatic -Description "Manage backend state on HAProxy"
Start-Service HAProxyDynAgent