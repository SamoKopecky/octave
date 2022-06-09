package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/octave/command/pause"
	"github.com/lukasl-dev/octave/command/play"
	"github.com/lukasl-dev/octave/command/resume"
	"github.com/lukasl-dev/octave/command/seek"
	"github.com/lukasl-dev/octave/command/stop"
	"github.com/lukasl-dev/octave/command/volume"
	"github.com/lukasl-dev/octave/config"
	"github.com/lukasl-dev/waterlink/v2"
)

// app is the main application struct. It holds all necessary dependencies to
// run the bot.
type app struct {
	// cfg is the application configuration.
	cfg config.Config

	// session is the discord session on which the bot is running on.
	session *discordgo.Session

	// conn is the active lavalink connection.
	conn *waterlink.Connection

	// client is the lavalink client.
	client *waterlink.Client

	// cmds is the commandHandler that is responsible for handling commands.
	cmds *commandHandler

	// sessionID is the current session's ID which is received on discordgo.Ready.
	sessionID string
}

// newApp returns a new app configured by the given config.Config.
func newApp(cfg config.Config) *app {
	return &app{cfg: cfg, cmds: newCommandHandler()}
}

func (a *app) run() (err error) {
	timeout, _ := time.ParseDuration("30s")
	err = waitForLavaLink(a.cfg, timeout)
	if err != nil {
		return err
	}
	a.session, err = discordgo.New(fmt.Sprintf("Bot %s", a.cfg.Token))
	if err != nil {
		return fmt.Errorf("failed to create discord session: %w", err)
	}

	a.registerHandlers()
	return a.session.Open()
}

// createClient tries to create a new waterlink.Client and defines it in the app.
func (a *app) createClient() (err error) {
	a.client, err = waterlink.NewClient(fmt.Sprintf("http://%s", a.cfg.Lavalink.Host), a.credentials())
	return err
}

// createConnection tries to create a new waterlink.Connection and defines it in
// the app.
func (a *app) createConnection() (err error) {
	a.conn, err = waterlink.Open(fmt.Sprintf("ws://%s", a.cfg.Lavalink.Host), a.credentials())
	return err
}

// credentials returns the waterlink.Credentials to use for client and connection.
func (a *app) credentials() waterlink.Credentials {
	return waterlink.Credentials{
		Authorization: a.cfg.Lavalink.Passphrase,
		UserID:        snowflake.MustParse(a.session.State.User.ID),
	}
}

// registerCommands registers all commands in the internal commandHandler.
func (a *app) registerCommands() {
	a.cmds.add(pause.Pause(pause.Deps{Conn: a.conn}))
	a.cmds.add(play.Play(play.Deps{Client: a.client, Conn: a.conn}))
	a.cmds.add(resume.Resume(resume.Deps{Conn: a.conn}))
	a.cmds.add(seek.Seek(seek.Deps{Conn: a.conn}))
	a.cmds.add(stop.Stop(stop.Deps{Conn: a.conn}))
	a.cmds.add(volume.Volume(volume.Deps{Conn: a.conn}))
}

func waitForLavaLink(cfg config.Config, maxTime time.Duration) error {
	for i := 0; i < int(maxTime.Seconds()); i++ {

		if isLavaLinkReady(cfg, time.Second) {
			return nil
		}
	}
	return errors.New("LavaLink is unreachable")
}

func isLavaLinkReady(cfg config.Config, timeout time.Duration) bool {
	host := strings.Split(cfg.Lavalink.Host, ":")[0]
	port := strings.Split(cfg.Lavalink.Host, ":")[1]
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		fmt.Println("Connecting error:", err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		fmt.Println("Found LavaLink at", net.JoinHostPort(host, port))
		return true
	}
	return false
}
