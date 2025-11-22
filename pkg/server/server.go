/*
 * Iptv-Proxy is a project to proxyfie an m3u file and to proxyfie an Xtream iptv service (client API).
 * Copyright (C) 2020  Pierre-Emmanuel Jacquier
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package server

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/jamesnetherton/m3u"
	"github.com/pierre-emmanuelJ/iptv-proxy/pkg/config"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

var defaultProxyfiedM3UPath = filepath.Join(os.TempDir(), uuid.NewV4().String()+".iptv-proxy.m3u")
var endpointAntiColision = strings.Split(uuid.NewV4().String(), "-")[0]

var (
	activeStreams     int64
	activeStreamsLock sync.Mutex
)

// Config represent the server configuration
type Config struct {
	*config.ProxyConfig

	// M3U service part
	playlist *m3u.Playlist
	// this variable is set only for m3u proxy endpoints
	track *m3u.Track
	// path to the proxyfied m3u file
	proxyfiedM3UPath string

	endpointAntiColision string

	channelRegistry        *channelRegistry
	httpClient             *http.Client
	hlsClient              *http.Client
	providerMaxConnections int
}

// NewServer initialize a new server configuration
func NewServer(config *config.ProxyConfig) (*Config, error) {
	var p m3u.Playlist
	if config.RemoteURL.String() != "" {
		var err error
		p, err = m3u.Parse(config.RemoteURL.String())
		if err != nil {
			return nil, err
		}
	}

	if trimmedCustomId := strings.Trim(config.CustomId, "/"); trimmedCustomId != "" {
		endpointAntiColision = trimmedCustomId
	}

	streamTransport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
		DisableKeepAlives:     false,
	}

	httpClient := &http.Client{
		Transport: streamTransport,
		Timeout:   0,
	}

	hlsTransport := &http.Transport{
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   20,
		MaxConnsPerHost:       0,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
		DisableKeepAlives:     false,
	}

	hlsClient := &http.Client{
		Transport: hlsTransport,
		Timeout:   0,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &Config{
		config,
		&p,
		nil,
		defaultProxyfiedM3UPath,
		endpointAntiColision,
		newChannelRegistry(),
		httpClient,
		hlsClient,
		0,
	}, nil
}
func streamStart(streamURL string, providerMaxConns int) {
	count := atomic.AddInt64(&activeStreams, 1)
	activeStreamsLock.Lock()
	defer activeStreamsLock.Unlock()

	log.Printf("[iptv-proxy] STREAM START: %s", streamURL)
	if providerMaxConns > 0 {
		log.Printf("[iptv-proxy] Active connections: %d/%d (%.0f%% capacity)",
			count, providerMaxConns, float64(count)/float64(providerMaxConns)*100)
		if count >= int64(providerMaxConns) {
			log.Printf("[iptv-proxy] WARNING: At or exceeding provider limit! Provider may reject connections.")
		}
	} else {
		log.Printf("[iptv-proxy] Active connections: %d (provider limit unknown)", count)
	}
}

func streamEnd(streamURL string, providerMaxConns int) {
	count := atomic.AddInt64(&activeStreams, -1)
	activeStreamsLock.Lock()
	defer activeStreamsLock.Unlock()

	log.Printf("[iptv-proxy] STREAM END: %s", streamURL)
	if providerMaxConns > 0 {
		log.Printf("[iptv-proxy] Active connections: %d/%d (%.0f%% capacity)",
			count, providerMaxConns, float64(count)/float64(providerMaxConns)*100)
	} else {
		log.Printf("[iptv-proxy] Active connections: %d", count)
	}
}

// Serve the iptv-proxy api
func (c *Config) Serve() error {
	if err := c.playlistInitialization(); err != nil {
		return err
	}

	router := gin.New()
	router.Use(gin.LoggerWithFormatter(channelAwareFormatter), gin.Recovery())
	router.Use(cors.Default())
	group := router.Group("/")
	c.routes(group)

	log.Printf("[iptv-proxy] Features:")
	log.Printf("[iptv-proxy]   - Channel-aware logging: enabled")
	log.Printf("[iptv-proxy]   - Channel registry: enabled")
	log.Printf("[iptv-proxy]   - Advanced error handling: enabled")
	if config.CacheFolder != "" {
		log.Printf("[iptv-proxy]   - Response caching: enabled")
	}
	if config.DebugLoggingEnabled {
		log.Printf("[iptv-proxy]   - Debug mode: enabled")
	}

	log.Printf("[iptv-proxy] Server is ready and listening on :%d", c.HostConfig.Port)

	return router.Run(fmt.Sprintf(":%d", c.HostConfig.Port))
}

func (c *Config) playlistInitialization() error {
	if len(c.playlist.Tracks) == 0 {
		return nil
	}

	f, err := os.Create(c.proxyfiedM3UPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.marshallInto(f, false)
}

// MarshallInto a *bufio.Writer a Playlist.
func (c *Config) marshallInto(into *os.File, xtream bool) error {
	filteredTrack := make([]m3u.Track, 0, len(c.playlist.Tracks))

	ret := 0
	into.WriteString("#EXTM3U\n") // nolint: errcheck
	for i, track := range c.playlist.Tracks {
		var buffer bytes.Buffer

		buffer.WriteString("#EXTINF:")                       // nolint: errcheck
		buffer.WriteString(fmt.Sprintf("%d ", track.Length)) // nolint: errcheck
		for i := range track.Tags {
			if i == len(track.Tags)-1 {
				buffer.WriteString(fmt.Sprintf("%s=%q", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
				continue
			}
			buffer.WriteString(fmt.Sprintf("%s=%q ", track.Tags[i].Name, track.Tags[i].Value)) // nolint: errcheck
		}

		uri, err := c.replaceURL(track.URI, i-ret, xtream)
		if err != nil {
			ret++
			log.Printf("ERROR: track: %s: %s", track.Name, err)
			continue
		}

		into.WriteString(fmt.Sprintf("%s, %s\n%s\n", buffer.String(), track.Name, uri)) // nolint: errcheck

		filteredTrack = append(filteredTrack, track)
	}
	c.playlist.Tracks = filteredTrack

	return into.Sync()
}

// ReplaceURL replace original playlist url by proxy url
func (c *Config) replaceURL(uri string, trackIndex int, xtream bool) (string, error) {
	oriURL, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	protocol := "http"
	if c.HTTPS {
		protocol = "https"
	}

	customEnd := strings.Trim(c.CustomEndpoint, "/")
	if customEnd != "" {
		customEnd = fmt.Sprintf("/%s", customEnd)
	}

	uriPath := oriURL.EscapedPath()
	if xtream {
		uriPath = strings.ReplaceAll(uriPath, c.XtreamUser.PathEscape(), c.User.PathEscape())
		uriPath = strings.ReplaceAll(uriPath, c.XtreamPassword.PathEscape(), c.Password.PathEscape())
	} else {
		uriPath = path.Join("/", c.endpointAntiColision, c.User.PathEscape(), c.Password.PathEscape(), fmt.Sprintf("%d", trackIndex), path.Base(uriPath))
	}

	basicAuth := oriURL.User.String()
	if basicAuth != "" {
		basicAuth += "@"
	}

	newURI := fmt.Sprintf(
		"%s://%s%s:%d%s%s",
		protocol,
		basicAuth,
		c.HostConfig.Hostname,
		c.AdvertisedPort,
		customEnd,
		uriPath,
	)

	newURL, err := url.Parse(newURI)
	if err != nil {
		return "", err
	}

	return newURL.String(), nil
}
