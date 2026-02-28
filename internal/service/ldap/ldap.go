package ldap

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"strings"

	"monolith/internal/config"

	ldapv3 "github.com/go-ldap/ldap/v3"
)

type Service struct {
	config config.LDAPConfig
}

type LDAPUser struct {
	Username string
	Email    string
	Name     string
}

func NewService(cfg config.LDAPConfig) *Service {
	return &Service{config: cfg}
}

func (s *Service) Enabled() bool {
	return s.config.Enabled
}

func (s *Service) AutoProvision() bool {
	return s.config.AutoProvision
}

func (s *Service) Authenticate(login, password string) (*LDAPUser, error) {
	conn, err := s.connect()
	if err != nil {
		return nil, fmt.Errorf("ldap connect: %w", err)
	}
	defer conn.Close()

	if err := conn.Bind(s.config.BindDN, s.config.BindPassword); err != nil {
		return nil, fmt.Errorf("ldap service bind: %w", err)
	}

	filter := strings.ReplaceAll(s.config.SearchFilter, "%s", ldapv3.EscapeFilter(login))
	searchRequest := ldapv3.NewSearchRequest(
		s.config.BaseDN,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases,
		0, 0, false,
		filter,
		[]string{"dn", s.config.UsernameAttribute, s.config.EmailAttribute, s.config.NameAttribute},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("ldap search: %w", err)
	}

	if len(sr.Entries) == 0 {
		return nil, ErrUserNotFound
	}

	if len(sr.Entries) > 1 {
		slog.Warn("LDAP search returned multiple entries, using first match", "login", login, "count", len(sr.Entries))
	}

	entry := sr.Entries[0]

	if err := conn.Bind(entry.DN, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return &LDAPUser{
		Username: entry.GetAttributeValue(s.config.UsernameAttribute),
		Email:    entry.GetAttributeValue(s.config.EmailAttribute),
		Name:     entry.GetAttributeValue(s.config.NameAttribute),
	}, nil
}

func (s *Service) connect() (*ldapv3.Conn, error) {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: s.config.SkipTLSVerify, //nolint:gosec // configurable for self-signed certs
		ServerName:         s.config.Host,
	}

	if s.config.Port == 636 {
		return ldapv3.DialTLS("tcp", addr, tlsConfig)
	}

	conn, err := ldapv3.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	if s.config.StartTLS {
		if err := conn.StartTLS(tlsConfig); err != nil {
			conn.Close()
			return nil, fmt.Errorf("ldap StartTLS: %w", err)
		}
	}

	return conn, nil
}
