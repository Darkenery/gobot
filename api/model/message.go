package model

type Message struct {
	MessageId             int              `json:"message_id"`
	From                  *User            `json:"from"`
	Date                  int64            `json:"date"`
	Chat                  *Chat            `json:"chat"`
	ForwardFrom           *User            `json:"forward_from"`
	ForwardFromChat       *Chat            `json:"forward_from_chat"`
	ForwardDate           int              `json:"forward_date"`
	ReplyToMessage        *Message         `json:"reply_to_message"`
	EditDate              int              `json:"edit_date"`
	Text                  string           `json:"text"`
	Entities              []*MessageEntity `json:"entities"`
	Audio                 Audio            `json:"audio"`
	Document              Document         `json:"document"`
	Photo                 []*PhotoSize     `json:"photo"`
	Sticker               *Sticker         `json:"sticker"`
	Video                 *Video           `json:"video"`
	Voice                 *Voice           `json:"voice"`
	Caption               string           `json:"caption"`
	Contact               *Contact         `json:"contact"`
	Location              *Location        `json:"location"`
	Venue                 *Venue           `json:"venue"`
	NewChatMember         *User            `json:"new_chat_member"`
	NewChatTitile         string           `json:"new_chat_titile"`
	NewChatPhoto          []*PhotoSize     `json:"new_chat_photo"`
	DeleteChatPhoto       bool             `json:"delete_chat_photo"`
	GroupChatCreated      bool             `json:"group_chat_created"`
	SupergroupChatCreated bool             `json:"supergroup_chat_created"`
	ChannelChatCreated    bool             `json:"channel_chat_created"`
	MigrateToChatId       int              `json:"migrate_to_chat_id"`
	MigrateFromChatId     int              `json:"migrate_from_chat_id"`
	PinnedMessage         *Message         `json:"pinned_message"`
}
