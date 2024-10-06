[Setup]
AppName=Koksmat
AppVersion=1.0.0
DefaultDirName={autopf}\koksmat
DefaultGroupName=Koksmat
OutputBaseFilename=koksmat-installer
OutputDir=.
Compression=lzma
SolidCompression=yes

[Files]
Source: "bin\koksmat.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\koksmat"; Filename: "{app}\koksmat.exe"

[Run]
Filename: "{app}\koksmat.exe"; Description: "Launch Koksmat"; Flags: nowait postinstall skipifsilent
