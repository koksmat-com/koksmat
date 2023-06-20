package powershell

import (
	"fmt"
	"strings"
)

type NewSharedMailboxResult struct {
	Name               string `json:"Name"`
	DisplayName        string `json:"DisplayName"`
	ExchangeObjectId   string `json:"ExchangeObjectId"`
	PrimarySmtpAddress string `json:"PrimarySmtpAddress"`
}
type EmptyResult struct {
}

type Member struct {
	User         string `json:"User"`
	AccessRights string `json:"AccessRights"`
	IsInherited  bool   `json:"IsInherited"`
}
type MembersResponse struct {
	Members []Member `json:"Members"`
}

type OwnersResponse struct {
	Owners string `json:"Owners"`
}

func PwshArray(members []string) string {

	pwshArray := strings.Join(members, ",")
	if (pwshArray) == "" {

		return "\"\""
	} else {
		return pwshArray
	}
}
func CreateSharedMailbox(Name string, DisplayName string, Alias string, Owners []string, Members []string, Readers []string) (result *NewSharedMailboxResult, err error) {
	powershellScript := "scripts/sharedmailboxes/create.ps1"
	powershellArguments := fmt.Sprintf(` -Name "%s" -DisplayName "%s"  -Alias "%s" -Members %s -Readers %s -Owners="%s"`, Name, DisplayName, Alias, PwshArray(Members), PwshArray(Readers), PwshArray(Owners))
	result, err = Run[NewSharedMailboxResult](powershellScript, powershellArguments)
	if err != nil {
		return result, err
	}

	return result, err
}

func DeleteSharedMailbox(ExchangeObjectId string) (result *EmptyResult, err error) {
	powershellScript := "scripts/sharedmailboxes/remove.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s`, ExchangeObjectId)
	result, err = Run[EmptyResult](powershellScript, powershellArguments)
	if err != nil {
		return result, err
	}

	return result, err
}

func UpdateSharedMailbox(ExchangeObjectId string, DisplayName string) (result *EmptyResult, err error) {
	powershellScript := "scripts/sharedmailboxes/update.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -DisplayName "%s"`, ExchangeObjectId, DisplayName)
	result, err = Run[EmptyResult](powershellScript, powershellArguments)
	if err != nil {
		return result, err
	}

	return result, err

}

func UpdateSharedMailboxPrimaryEmailAddress(ExchangeObjectId string, Email string) (result *EmptyResult, err error) {
	powershellScript := "scripts/sharedmailboxes/updateprimaryemail.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Email "%s"`, ExchangeObjectId, Email)
	result, err = Run[EmptyResult](powershellScript, powershellArguments)
	if err != nil {
		return result, err
	}

	return result, err

}

func AddSharedMailboxMembers(ExchangeObjectId string, Members []string) (members *MembersResponse, err error) {
	powershellScript := "scripts/sharedmailboxes/addmembers.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Members %s`, ExchangeObjectId, PwshArray(Members))
	members, err = Run[MembersResponse](powershellScript, powershellArguments)
	return members, err
}

func AddSharedMailboxReaders(ExchangeObjectId string, Readers []string) (members *MembersResponse, err error) {
	powershellScript := "scripts/sharedmailboxes/addreaders.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Readers %s`, ExchangeObjectId, PwshArray(Readers))
	members, err = Run[MembersResponse](powershellScript, powershellArguments)
	return members, err
}

func SetSharedMailboxOwners(ExchangeObjectId string, Owners []string) (res *OwnersResponse, err error) {
	powershellScript := "scripts/sharedmailboxes/addowners.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Owners %s`, ExchangeObjectId, PwshArray(Owners))
	res, err = Run[OwnersResponse](powershellScript, powershellArguments)
	return res, err
}

func RemoveSharedMailboxMembers(ExchangeObjectId string, Members []string) (members *MembersResponse, err error) {
	powershellScript := "scripts/sharedmailboxes/removemembers.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Members %s`, ExchangeObjectId, PwshArray(Members))
	members, err = Run[MembersResponse](powershellScript, powershellArguments)
	return members, err
}

func RemoveSharedMailboxReaders(ExchangeObjectId string, Readers []string) (members *MembersResponse, err error) {
	powershellScript := "scripts/sharedmailboxes/removereaders.ps1"
	powershellArguments := fmt.Sprintf(` -ExchangeObjectId %s -Readers %s`, ExchangeObjectId, PwshArray(Readers))
	members, err = Run[MembersResponse](powershellScript, powershellArguments)
	return members, err
}
