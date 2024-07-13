

$EXCHAPPID = $env:EXCHAPPID
$EXCHORGANIZATION = $env:EXCHORGANIZATION
$EXCHCERTIFICATEPASSWORD = $env:EXCHCERTIFICATEPASSWORD
$EXCHCERTIFICATEPATH = "$PSScriptRoot/certificate.pfx"
$bytes = [Convert]::FromBase64String($ENV:EXCHCERTIFICATE)
[IO.File]::WriteAllBytes($EXCHCERTIFICATEPATH, $bytes)

Write-Output "Connecting to Exchange for $EXCHORGANIZATION"

if (($EXCHCERTIFICATEPASSWORD -ne $null) -and ($EXCHCERTIFICATEPASSWORD -ne "") ){
	Connect-ExchangeOnline -CertificateFilePath $EXCHCERTIFICATEPATH  -AppID $EXCHAPPID -Organization $EXCHORGANIZATION -ShowBanner:$false -CertificatePassword (ConvertTo-SecureString -String $EXCHCERTIFICATEPASSWORD -AsPlainText -Force)
}else{
	Connect-ExchangeOnline -CertificateFilePath $EXCHCERTIFICATEPATH  -AppID $EXCHAPPID -Organization $EXCHORGANIZATION -ShowBanner:$false #   -BypassMailboxAnchoring:$true

}							
				