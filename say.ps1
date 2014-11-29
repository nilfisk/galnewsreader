param(
	[Parameter(Mandatory=$true)]
	$Text
)

$object = New-Object -ComObject SAPI.SpVoice
$object.Speak($text) | Out-Null
