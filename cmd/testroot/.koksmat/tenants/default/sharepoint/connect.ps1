

$PNPAPPID = $env:PNPAPPID
$PNPTENANTID = $env:PNPTENANTID
$PNPCERTIFICATEPATH = "$($PSScriptRoot)/pnp.pfx"
$PNPSITE = $env:PNPSITE
$bytes = [Convert]::FromBase64String($ENV:PNPCERTIFICATE)
[IO.File]::WriteAllBytes($PNPCERTIFICATEPATH, $bytes)

write-output "Connecting to $PNPSITE"
Connect-PnPOnline -Url $PNPSITE  -ClientId $PNPAPPID -Tenant $PNPTENANTID -CertificatePath "$PNPCERTIFICATEPATH"
			
							
				