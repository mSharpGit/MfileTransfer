;; RSA Token Automation.ahk
;; # autohotkey.exe
;; # Ref: http://www.autohotkey.com/board/topic/59612-simple-debug-console-output/
DebugMessage(str)
{
 global h_stdout
 DebugConsoleInitialize()  ; start console window if not yet started
 str .= "`n" ; add line feed
 DllCall("WriteFile", "uint", h_Stdout, "uint", &str, "uint", StrLen(str), "uint*", BytesWritten, "uint", NULL) ; write into the console
 WinSet, Bottom,, ahk_id %h_stout%  ; keep console on bottom
}

DebugConsoleInitialize()
{
   global h_Stdout	 ; Handle for console
   static is_open = 0  ; toogle whether opened before
   if (is_open = 1)	 ; yes, so don't open again
	 return
	 
   is_open := 1	
   ; two calls to open, no error check (it's debug, so you know what you are doing)
   DllCall("AttachConsole", int, -1, int)
   DllCall("AllocConsole", int)

   dllcall("SetConsoleTitle", "str","Paddy Debug Console")	; Set the name. Example. Probably could use a_scriptname here 
   h_Stdout := DllCall("GetStdHandle", "int", -11) ; get the handle
   WinSet, Bottom,, ahk_id %h_stout%	  ; make sure it's on the bottom
   WinActivate,Lightroom   ; Application specific; I need to make sure this application is running in the foreground. YMMV
   return
}

;; RSA
;;file := FileOpen(tt.txt, "w")
Run, "C:\Program Files (x86)\RSA SecurID Software Token\SecurID.exe"

WinWait, 000140089513 - RSA SecurID Token, 
IfWinNotActive, 000140089513 - RSA SecurID Token, 
WinActivate, 000140089513 - RSA SecurID Token, 
WinWaitActive, 000140089513 - RSA SecurID Token, 
;;MsgBox, "sleeping"
;;Sleep, 3000
Send, 1926
;;MsgBox, "sleeping"
Sleep, 100
;;MsgBox, "awake"
Send, {ENTER}
Sleep, 500
;;Send, {CTRLDOWN}c{CTRLUP}

Send, ^c
Sleep, 1000


;;DebugMessage(X)
Send, {ALTDOWN}{F4}{ALTUP}

Sleep, 1000

;;Send, ^v

X := clipboard
;;OutputDebug, X
Sleep 100
;;MsgBox, %X%
Sleep 100
IniWrite, %X%, C:\Users\masterserver\Desktop\master-files-automation\hotkey\token.ini, token, value
Sleep 100

;;FileName := C:\Users\masterserver\Desktop\master-automation\tot.txt
;;FileSelectFile, %FileName% , S16,, Create a new file:
;;if (%FileName% = "")
;;    return
;;file := FileOpen( tot.txt, "w")
;;if !IsObject(file)
;;{
 ;;   MsgBox Can't open "tot.txt" for writing.
 ;;   return
;;}
;;TestString := %X% ;;"This is a test string.`r`n"  ; When writing a file this way, use `r`n rather than `n to start a new line.
;;file.Write(TestString)
;;file.Close()
;;EnvSet, TOKEN, %X%



;;DebugMessage(X)


;;file.write(clipboard)
;;Run, CMD.exe 
;;WinWait, ahk_exe cmd.exe
;;Send, cd C:\Users\masterserver\Desktop\master-automation{enter}
;;Send, sftp.bat ^v MCI.AR.T461.M.E0073275.D181017* /0073275/production/download/T461 > tt.txt
;;sleep 500
;;Send, {enter}
;;Send, copy MCI.AR.T461.M.E0073275.D181017* TT461
;;Send, {enter}
;;Send, {ALTDOWN}{F4}{ALTUP}
return


