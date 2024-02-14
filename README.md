# Skylab Misfin Client

To read about Skylab, see the page on my Gemini capsule:
=> gemini://satch.xyz/skylab/ Skylab on Gemini

=> gemini://bbs.geminispace.org/s/Skylab Issue Tracker

## Packages, functions, types
- cgi endpoints
  - archive
  - compose
  - delete
  - drafts
  - home
  - inbox
  - mail
  - reply
  - send
  - sent
  - settings
  - tag
  - tags
  - trash
  - unread
- fetch
  - type Message struct
  - func FetchGembox
  - func isNewMessage
  - func ParseGembox
  - func isolateAddress
- config
  - type Config struct
  - func ExportConfig
- send
  - func SendSkymail
  - func GetAddressFromCert
  - func SkymailToString
  - func ParseComposePath
- tools
  - type Skymail struct
  - func GenerateFilename
  - func GenerateID
  - func ParseSkymail
  - func IsolateAddress
  - func PrintLinkAndPreview
  - func GetMessageLinkLine
  - func GetMailPreview
  - func Tag
  - func Untag
  - func isDirEmpty
  - func HasTag
  - func WalkTag
  - func FormatSender
  - func CountTags
  - func UntagAll
  - func ListTags
  - type UserHashes
  - func Authenticate
  - func AddUserHash

Next steps:
- Better new user creation
- Better message fetching options IMAP-like protocol, general compatibility with other server implementations
- more config options (custom paths)
