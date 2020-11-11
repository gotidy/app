package config

import (
	"net"
	"net/url"
	"strconv"
)

// NetAddress
type NetAddress struct {
	Host string
	Port int
}

func (a NetAddress) Network() string {
	return ""
}

// String combines host and port into a network address of the
// form "host:port". If host contains a colon, as found in literal
// IPv6 addresses, then JoinHostPort returns "[host]:port".
func (a NetAddress) String() string {
	s := url.PathEscape(a.Host)
	if a.Port != 0 {
		s = net.JoinHostPort(s, strconv.Itoa(a.Port))
	}
	return s
}

type UserCredential struct {
	User     string
	Password *string
}

func NewUserCredential(user string, password string) UserCredential {
	return UserCredential{
		User:     user,
		Password: &password,
	}
}

// String combines user and password into the form "user:password".
func (u UserCredential) String() string {
	if u.Password != nil {
		return url.UserPassword(u.User, *u.Password).String()
	}
	return url.User(u.User).String()
}

// Connection's address and user credential
type Connection struct {
	User    UserCredential
	Address NetAddress
}

// String combines user and address into the form "user:password@host:port".
func (c Connection) String() string {
	s := c.Address.String()
	if u := c.User.String(); u != "" {
		s = u + "@" + s
	}
	return s
}

// String combines scheme, user and address into the form "user:password@host:port".
func (c Connection) URL(scheme string) string {
	return scheme + "://" + c.String() + "/"
}
