package server

import (
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/jamesnetherton/m3u"
)

const channelNameContextKey = "channel_name"

type channelRegistry struct {
	mu    sync.RWMutex
	names map[string]string
}

func newChannelRegistry() *channelRegistry {
	return &channelRegistry{
		names: make(map[string]string),
	}
}

func (r *channelRegistry) remember(uri, name string) {
	if r == nil || name == "" || uri == "" {
		return
	}

	ids := extractStreamIdentifiers(uri)
	if len(ids) == 0 {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		if id == "" {
			continue
		}
		r.names[id] = name
	}
}

func (r *channelRegistry) rememberIdentifiers(name string, ids ...string) {
	if r == nil || name == "" {
		return
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range ids {
		clean := sanitizeIdentifier(id)
		if clean == "" {
			continue
		}
		r.names[clean] = name
	}
}

func (r *channelRegistry) lookup(id string) string {
	if r == nil {
		return ""
	}

	key := sanitizeIdentifier(id)
	if key == "" {
		return ""
	}

	r.mu.RLock()
	name := r.names[key]
	r.mu.RUnlock()
	if name != "" {
		return name
	}

	if dot := strings.Index(key, "."); dot > 0 {
		trimmed := key[:dot]
		r.mu.RLock()
		name = r.names[trimmed]
		r.mu.RUnlock()
	}

	return name
}

func (c *Config) rememberPlaylistChannels(pl *m3u.Playlist) {
	if c == nil || pl == nil {
		return
	}

	for i := range pl.Tracks {
		c.rememberTrackChannel(&pl.Tracks[i])
	}
}

func (c *Config) rememberTrackChannel(track *m3u.Track) {
	if c == nil || c.channelRegistry == nil || track == nil {
		return
	}

	c.channelRegistry.remember(track.URI, track.Name)
}

func (c *Config) rememberIdentifiers(name string, ids ...string) {
	if c == nil || c.channelRegistry == nil {
		return
	}

	c.channelRegistry.rememberIdentifiers(name, ids...)
}

func (c *Config) annotateChannel(ctx *gin.Context, ids ...string) string {
	if c == nil || c.channelRegistry == nil || ctx == nil {
		return ""
	}

	for _, id := range ids {
		if id == "" {
			continue
		}

		if name := c.channelRegistry.lookup(id); name != "" {
			ctx.Set(channelNameContextKey, name)
			return name
		}
	}

	return ""
}

func extractStreamIdentifiers(raw string) []string {
	u, err := url.Parse(raw)
	if err != nil {
		return nil
	}

	base := sanitizeIdentifier(path.Base(u.Path))
	if base == "" {
		return nil
	}

	identifiers := []string{base}
	if dot := strings.Index(base, "."); dot > 0 {
		identifiers = append(identifiers, base[:dot])
	}

	if rawID := sanitizeIdentifier(path.Base(raw)); rawID != "" && rawID != base {
		identifiers = append(identifiers, rawID)
	}

	return identifiers
}

func sanitizeIdentifier(id string) string {
	return strings.Trim(strings.TrimSpace(id), "/")
}
