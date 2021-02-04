package actions

import (
	"errors"
	"raiki"
)

var (
	ErrInvalidJson            = errors.New("invalid json")
	ErrInsufficientPermission = errors.New("insufficient permission")

	ErrBanDelete  = errors.New("could not delete ban")
	ErrBanGet     = errors.New("could not get bans")
	ErrMuteDelete = errors.New("could not delete mute")
	ErrMuteGet    = errors.New("could not get mutes")

	ErrHistoryGet = errors.New("could not get history")

	ErrChatSave  = errors.New("could not save chat")
	ErrChatShort = errors.New("chat too short")

	ErrCommunityExists         = errors.New("community exists")
	ErrCommunityLimit          = errors.New("too many communities")
	ErrCommunitySave           = errors.New("could not save community")
	ErrCommunityNonExist       = errors.New("community doesn't exist")
	ErrCommunityName           = errors.New("community name invalid")
	ErrCommunitySlug           = errors.New("community slug invalid")
	ErrCommunityDescription    = errors.New("community description invalid")
	ErrCommunityWelcomeMessage = errors.New("community welcome message invalid")
	ErrCommunityState          = errors.New("could not retrive community state")
	ErrCommunityList           = errors.New("could not retrieve community list")

	ErrCommunityMedia = errors.New("no media currently playing")

	ErrCommunityWaitlistJoinPlaylist      = errors.New("must have at least 1 playlist")
	ErrCommunityWaitlistJoinPlaylistItems = errors.New("must have at least 1 item in active playlist")

	ErrStaffSave = errors.New("could not save staff")
	ErrStaffGet  = errors.New("could not get staff")

	ErrUserNonExist = errors.New("user doesn't exist")

	ErrPlaylistActive     = errors.New("playlist already active")
	ErrPlaylistGet        = errors.New("could not get playlist")
	ErrPlaylistSave       = errors.New("could not save playlist")
	ErrPlaylistsGet       = errors.New("could not get playlists")
	ErrPlaylistName       = errors.New("playlist name invalid")
	ErrPlaylistItemLimit  = errors.New("too many items")
	ErrPlaylistLimit      = errors.New("too many playlists")
	ErrPlaylistDelete     = errors.New("could not delete playlist")
	ErrPlaylistDeleteLast = errors.New("you cannot delete your last playlist")
	ErrPlaylistNonExist   = errors.New("playlist doesn't exist")
	ErrPlaylistOwner      = errors.New("you do not own this playlist")

	ErrPlaylistItemNonExist        = errors.New("playlist item doesn't exist")
	ErrPlaylistItemExists          = errors.New("playlist item already exists")
	ErrPlaylistItemInvalidPosition = errors.New("invlid position")

	ErrMediaType     = errors.New("invalid media type")
	ErrMediaDownload = errors.New("could not retrieve media data")
	ErrMediaSave     = errors.New("could not save media")
	ErrMediaSearch   = errors.New("could not search for media")

	ErrBanned = errors.New("banned")
	ErrMuted  = errors.New("muted")

	SocketMap = map[string]raiki.ActionFunc{
		"echo":                            Echo,
		"chat.send":                       ChatSend,
		"community.create":                CommunityCreate,
		"community.get.bans":              CommunityGetBans,
		"community.get":                   CommunityGet,
		"community.get.history":           CommunityGetHistory,
		"community.get.media":             CommunityGetMedia,
		"community.get.mutes":             CommunityGetMutes,
		"community.get.staff":             CommunityGetStaff,
		"community.get.users":             CommunityGetUsers,
		"community.join":                  CommunityJoin,
		"community.list":                  CommunityList,
		"community.update.description":    CommunityUpdateDescription,
		"community.update.djCycle":        CommunityUpdateDjCycle,
		"community.update.name":           CommunityUpdateName,
		"community.update.nsfw":           CommunityUpdateNsfw,
		"community.update.waitlistLocked": CommunityUpdateWaitlistLocked,
		"community.update.welcomeMessage": CommunityUpdateWelcomeMessage,
		"community.vote.meh":              CommunityVoteMeh,
		"community.vote.save":             CommunityVoteSave,
		"community.vote.woot":             CommunityVoteWoot,
		"community.waitlist.join":         CommunityWaitlistJoin,
		"community.waitlist.leave":        CommunityWaitlistLeave,
		"media.search":                    MediaSearch,
		"playlist.get":                    PlaylistGet,
		"playlist.activate":               PlaylistActivate,
		"playlist.create":                 PlaylistCreate,
		"playlist.delete":                 PlaylistDelete,
		"playlist.insert":                 PlaylistInsert,
		"playlist.item.delete":            PlaylistItemDelete,
		"playlist.item.edit":              PlaylistItemEdit,
		"playlist.item.move":              PlaylistItemMove,
		"self.get":                        SelfGet,
		"user.get":                        UserGet,
	}
)
