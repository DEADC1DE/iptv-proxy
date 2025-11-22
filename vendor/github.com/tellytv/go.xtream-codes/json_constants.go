package xtreamcodes

const (
	StructFields = "Fields"
)

// ServerInfoFields contains constants for ServerInfo struct JSON field names
const (
	ServerInfoFieldHTTPSPort    = "https_port"
	ServerInfoFieldPort         = "port"
	ServerInfoFieldProcess      = "process"
	ServerInfoFieldRTMPPort     = "rtmp_port"
	ServerInfoFieldProtocol     = "server_protocol"
	ServerInfoFieldTimeNow      = "time_now"
	ServerInfoFieldTimestampNow = "timestamp_now"
	ServerInfoFieldTimezone     = "timezone"
	ServerInfoFieldURL          = "url"
)

// UserInfoFields contains constants for UserInfo struct JSON field names
const (
	UserInfoFieldActiveConnections   = "active_cons"
	UserInfoFieldAllowedOutputFormats = "allowed_output_formats"
	UserInfoFieldAuth                = "auth"
	UserInfoFieldCreatedAt           = "created_at"
	UserInfoFieldExpDate             = "exp_date"
	UserInfoFieldIsTrial             = "is_trial"
	UserInfoFieldMaxConnections      = "max_connections"
	UserInfoFieldMessage             = "message"
	UserInfoFieldPassword            = "password"
	UserInfoFieldStatus              = "status"
	UserInfoFieldUsername            = "username"
)

// StreamFields contains constants for Stream struct JSON field names
const (
	StreamFieldAdded            = "added"
	StreamFieldCategoryID       = "category_id"
	StreamFieldCategoryIDs      = "category_ids"
	StreamFieldContainerExtension = "container_extension"
	StreamFieldCustomSID        = "custom_sid"
	StreamFieldDirectSource     = "direct_source"
	StreamFieldEPGChannelID     = "epg_channel_id"
	StreamFieldIcon             = "stream_icon"
	StreamFieldID               = "stream_id"
	StreamFieldName             = "name"
	StreamFieldNumber           = "num"
	StreamFieldRating           = "rating"
	StreamFieldRating5Based     = "rating_5based"
	StreamFieldStreamType       = "stream_type"
	StreamFieldTVArchive        = "tv_archive"
	StreamFieldTVArchiveDuration = "tv_archive_duration"
)

// CategoryFields contains constants for Category struct JSON field names
const (
	CategoryFieldCategoryID   = "category_id"
	CategoryFieldCategoryName = "category_name"
	CategoryFieldParentID     = "parent_id"
)

// SeriesInfoFields contains constants for SeriesInfo struct JSON field names
const (
	SeriesInfoFieldBackdropPath    = "backdrop_path"
	SeriesInfoFieldCast            = "cast"
	SeriesInfoFieldCategoryID      = "category_id"
	SeriesInfoFieldCover           = "cover"
	SeriesInfoFieldDirector        = "director"
	SeriesInfoFieldEpisodes        = "episodes"
	SeriesInfoFieldGenre           = "genre"
	SeriesInfoFieldLast_Modified   = "last_modified"
	SeriesInfoFieldName            = "name"
	SeriesInfoFieldPlot            = "plot"
	SeriesInfoFieldRating          = "rating"
	SeriesInfoFieldRating5Based    = "rating_5based"
	SeriesInfoFieldReleaseDate     = "releaseDate"
	SeriesInfoFieldYouTubeTrailer  = "youtube_trailer"
)

// SeriesEpisodeFields contains constants for SeriesEpisode struct JSON field names
const (
	SeriesEpisodeFieldAdded           = "added"
	SeriesEpisodeFieldBitrate         = "bitrate"
	SeriesEpisodeFieldContainerExtension = "container_extension"
	SeriesEpisodeFieldCustomSID       = "custom_sid"
	SeriesEpisodeFieldDirectSource    = "direct_source"
	SeriesEpisodeFieldDuration        = "duration"
	SeriesEpisodeFieldDurationSecs    = "duration_secs"
	SeriesEpisodeFieldEpisodeNumber   = "episode_num"
	SeriesEpisodeFieldID              = "id"
	SeriesEpisodeFieldInfo            = "info"
	SeriesEpisodeFieldSeason          = "season"
	SeriesEpisodeFieldTitle           = "title"
)

// VideoOnDemandInfoFields contains constants for VideoOnDemandInfo struct JSON field names
const (
	VideoOnDemandInfoFieldInfo  = "info"
	VideoOnDemandInfoFieldMovie_Data = "movie_data"
)

// SeriesFields contains constants for Series struct JSON field names
const (
	SeriesFieldCategoryID   = "category_id"
	SeriesFieldCategoryIDs  = "category_ids"
	SeriesFieldCover        = "cover"
	SeriesFieldEpisodeRunTime = "episode_run_time"
	SeriesFieldID           = "series_id"
	SeriesFieldLast_Modified = "last_modified"
	SeriesFieldName         = "name"
	SeriesFieldNumber       = "num"
	SeriesFieldPlot         = "plot"
	SeriesFieldRating       = "rating"
	SeriesFieldRating5Based = "rating_5based"
	SeriesFieldYouTubeTrailer = "youtube_trailer"
)
