package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	domains := make(DomainStat)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User

		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		hostname := strings.ToLower(strings.Split(user.Email, "@")[1])
		dom := strings.ToLower(strings.Split(hostname, ".")[1])

		if dom == domain {
			domains[hostname]++
		}
	}

	return domains, nil
}
