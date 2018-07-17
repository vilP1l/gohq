package gohq

import (
	"net/url"
)

var (
	EndpointBase = "https://api-quiz.hype.space/"

	EndpointUsers         = EndpointBase + "users/"
	EndpointMe            = EndpointUsers + "me/"
	EndpointPayouts       = EndpointMe + "payouts/"
	EndpointShows         = EndpointBase + "shows/"
	EndpointSchedule      = EndpointShows + "now?type=hq"
	EndpointAvatarURL     = EndpointMe + "avatarUrl/"
	EndpointFriends       = EndpointBase + "friends/"
	EndpointVerifications = EndpointBase + "verifications/"
	EndpointEasterEggs    = EndpointBase + "easter-eggs/"
	EndpointAWS           = EndpointBase + "credentials/s3"
	EndpointMakeItRain    = EndpointEasterEggs + "makeItRain/"
	EndpointTokens        = EndpointBase + "tokens/"
	EndpointToken         = EndpointMe + "token/"
	EndpointGifts         = EndpointBase + "gifts/"
	EndpointDrops         = EndpointGifts + "drops/"
	EndpointDevices = EndpointMe + "devices/"

	EndpointUser          = func(uID string) string { return EndpointUsers + uID + "/" }
	EndpointClaim         = func(gID string) string { return EndpointDrops + gID + "/claims" }
	EndpointFriend        = func(uID string) string { return EndpointFriends + uID + "/" }
	EndpointDocuments = func(uID string) string {return EndpointUsers + uID + "/payouts/documents"}
	EndpointFriendRequest = func(uID string) string { return EndpointFriend(uID) + "requests/" }
	EndpointSearchUser    = func(query string) string { return EndpointUsers + "?q=" + url.QueryEscape(query) }
)
