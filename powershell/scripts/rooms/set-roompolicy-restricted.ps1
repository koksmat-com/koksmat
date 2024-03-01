param (
    [Parameter(Mandatory = $true)]
    [string]$Mail,
    [Parameter(Mandatory = $true)]
    [string[]]$RestrictedTo,
    [Parameter(Mandatory = $true)]
    [string]$MailTip,
    [Parameter(Mandatory = $true)]
    [int]$BookingWindowInDays
)


write-host "Processing" $Mail
Set-Mailbox $mail  -MailTip $MailTip
Set-CalendarProcessing $mail  -DeleteComments $false -AutomateProcessing AutoAccept -AllRequestInPolicy $false  -AllBookInPolicy $false -BookInPolicy $restrictedTo -BookingWindowInDays 601 -ResourceDelegates $null       
  
