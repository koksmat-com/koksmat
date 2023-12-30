
param (

    [Parameter(Mandatory = $true)]
    [string]$kitchen,
    [Parameter(Mandatory = $true)]
    [string]$SiteURL,
    [Parameter(Mandatory = $true)]
    [string]$tenantDomain
)
$ErrorActionPreference = "Stop"
$kitchenRoot = $env:KITCHENROOT
set-location "$kitchenRoot/$kitchen"
$location = get-location 
. $location/.sharepoint/tenants/$tenantDomain/env.ps1
Connect-PnPOnline -Url $SiteURL  -ClientId $PNPAPPID -Tenant $PNPTENANTID -CertificatePath "$PNPCERTIFICATEPATH"

$site = Get-PnPSite -Includes RootWeb,ServerRelativeUrl,GroupId,HubSiteId,IsHubSite,SensitivityLabelInfo,SecondaryContact,Owner
$info  = @{
    webUrl = $site.Url
    Title = $site.RootWeb.Title
    
}

Invoke-PnPSiteTemplate -Path "$psscriptroot/templates/sharepoint/intra365/template.xml" -excludeHandlers SiteFooter,ApplicationLifecycleManagement,Navigation
convertto-json -InputObject $true
