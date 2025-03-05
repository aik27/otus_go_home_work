package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	i := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}
		result[i] = user
		i++
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, err
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched := strings.Contains(user.Email, domain)

		if matched {
			email := strings.SplitN(user.Email, "@", 2)
			if len(email) >= 2 {
				num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
				num++
				result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
			}
		}
	}
	return result, nil
}
