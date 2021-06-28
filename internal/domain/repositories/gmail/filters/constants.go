package filters

// Gmail search filters

const (
	LabelRead   = "Label:read"      // Readed mails
	LabelUnread = "Label:unread"    // Unreaded mails
	AfterDate   = "After:"          // After yyyy/mm/dd
	File        = "Filename: "      // extension or name of a specific one
	HasAttach   = "Has:attachment " // Attachments Specific descriptor
	From        = "From:"           // Sender
	To          = "To:"             // Receiver
)
