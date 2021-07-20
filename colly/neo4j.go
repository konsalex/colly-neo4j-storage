package neo4j

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Storage implements a Neo4j storage backend for colly
type Storage struct {
	// Database name, defaults to "neo4j"
	Database string
	// URI, with protocol included, ex. bolt://localhost:7687
	URI string
	// Username of the user
	Username string
	// Password of the user
	Password string
	// Neo4j Driver pointer
	driver *neo4j.Driver
}

// Init initializes the Neo4j driver
func (s *Storage) Init() error {

	var err error

	driver, err := neo4j.NewDriver(s.URI, neo4j.BasicAuth(s.Username, s.Password, ""))

	if err != nil {
		fmt.Println(err)
		return err
	}

	if s.Database == "" {
		fmt.Println("Database is not defined, using: neo4j")
		s.Database = "neo4j"
	}

	s.driver = &driver

	return err
}

// Visited implements colly/storage.Visited()
func (s *Storage) Visited(requestID uint64) error {

	driver := *s.driver

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: s.Database})
	defer session.Close()

	createVisitedQuery := `
	MERGE (a:Visited { requestID: $requestID, visited: $visited })
	`

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(createVisitedQuery, map[string]interface{}{
			"requestID": strconv.FormatUint(requestID, 10),
			"visited":   true,
		})
		return nil, err
	})

	if err != nil {
		fmt.Println(err)
		return err
	}

	return err
}

// IsVisited implements colly/storage.IsVisited()
func (s *Storage) IsVisited(requestID uint64) (bool, error) {

	driver := *s.driver

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: s.Database})
	defer session.Close()

	stringRequestID := strconv.FormatUint(requestID, 10)

	visitedQuery := `
	MATCH (a:Visited { requestID: $requestID})
	RETURN a.visited as visited`

	result, err := session.Run(visitedQuery, map[string]interface{}{
		"requestID": stringRequestID,
	})

	if err != nil {
		return false, err
	}

	if result.Next() {
		visited, ok := result.Record().Get("visited") // a.visited
		if ok {
			return visited.(bool), nil
		}
		return false, errors.New("result was not ok")
	}

	return false, nil
}

// Cookies implements colly/storage.Cookies()
func (s *Storage) Cookies(u *url.URL) string {

	driver := *s.driver

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: s.Database})
	defer session.Close()

	cookieQuery := `
	MATCH (a:Cookies { host: $host})
	RETURN a.cookies as cookies LIMIT 1`

	result, err := session.Run(cookieQuery, map[string]interface{}{
		"host": u.Host,
	})

	if err != nil {
		fmt.Println("Error in Cookie Retrieval")
		fmt.Println(err)
		return ""
	}

	if result.Next() {
		cookies, ok := result.Record().Get("cookies")
		if ok {
			return cookies.(string)
		}
		fmt.Println("Result was not ok")
	}

	return ""
}

// SetCookies implements colly/storage.SetCookies()
func (s *Storage) SetCookies(u *url.URL, cookies string) {

	driver := *s.driver

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: s.Database})
	defer session.Close()

	createVisitedQuery := `
	MERGE (a:Cookies { host: $host, cookies: $cookies })
	`

	_, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		_, err := tx.Run(createVisitedQuery, map[string]interface{}{
			"host":    u.Host,
			"cookies": cookies,
		})
		return nil, err
	})

	if err != nil {
		fmt.Println(err)
	}
}
