package bot

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

type BotDatabase struct {
	mutex sync.Mutex
	db    *sql.DB
}

type AccessRecord struct {
	ID            int
	UserID        string
	UserName      string
	UserAvatarUrl string
	AccessKey     string
	ScopeNews     sql.NullTime
	ScopeZerotier sql.NullTime
}

type ZeroTierNetworkRecord struct {
	ID          int
	NetworkID   string
	NetworkName string
}

type ZeroTierMemberRecord struct {
	ID         int
	NetworkID  string
	MemberID   string
	MemberName string
	UserID     string
}

func (r *AccessRecord) IsNewsScope() bool {
	return r.ScopeNews.Valid
}

func (r *AccessRecord) IsZerotierScope() bool {
	return r.ScopeZerotier.Valid
}

const create string = `
	CREATE TABLE IF NOT EXISTS access (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT NOT NULL,
		user_name TEXT NOT NULL,
		avatar_url TEXT NOT NULL,
		access_key TEXT NOT NULL,
		scope_news DATETIME DEFAULT NULL,
		scope_zerotier DATETIME DEFAULT NULL
	);

	CREATE TABLE IF NOT EXISTS network (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		network_id TEXT NOT NULL,
		network_name TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS member (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		network_id TEXT NOT NULL,
		member_id TEXT NOT NULL,
		member_name TEXT NOT NULL,
		user_id TEXT NOT NULL
	);
`

// Scope constants
type Scope int8

const (
	SCOPE_NEWS Scope = iota
	SCOPE_ZEROTIER
)

func OpenAccessDb() (BotDatabase, error) {

	db, err := sql.Open("sqlite3", "/app/data/tony.db")
	// db, err := sql.Open("sqlite3", "../database/tony.db")
	if err != nil {
		return BotDatabase{}, err
	}

	if _, err := db.Exec(create); err != nil {
		return BotDatabase{}, err
	}

	return BotDatabase{
		db: db,
	}, nil
}

func (db *BotDatabase) Close() error {
	return db.db.Close()
}

// Validate Access Key and get the username and scope
func (db *BotDatabase) ValidateAccessKey(accessKey string) (AccessRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var record AccessRecord
	record.AccessKey = accessKey

	err := db.db.QueryRow("SELECT id, user_id, user_name, scope_news, scope_zerotier FROM access WHERE access_key = ?", accessKey).Scan(&record.ID, &record.UserID, &record.UserName, &record.ScopeNews, &record.ScopeZerotier)
	if err != nil {
		return record, err
	}

	return record, nil
}

// Validate Access Key and get the username and scope
func (db *BotDatabase) GetUserAccess(userId string) (AccessRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var record AccessRecord
	record.UserID = userId

	err := db.db.QueryRow("SELECT id, user_name, avatar_url, scope_news, scope_zerotier FROM access WHERE user_id = ?", userId).Scan(&record.ID, &record.UserName, &record.UserAvatarUrl, &record.ScopeNews, &record.ScopeZerotier)
	if err != nil {
		return record, err
	}

	return record, nil
}

// Add User and Access Key to the database
func (db *BotDatabase) AddAccess(userId string, userName string, scope Scope, userAvatarUrl string) (AccessRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Check if user already exists
	var count int
	err := db.db.QueryRow("SELECT COUNT(*) FROM access WHERE user_id = ?", userId).Scan(&count)
	if err != nil {
		return AccessRecord{}, err
	}

	var access = AccessRecord{
		UserID:        userId,
		UserName:      userName,
		UserAvatarUrl: userAvatarUrl,
	}

	// If user exists, update the access key
	if count > 0 {
		switch scope {
		case SCOPE_NEWS:
			_, err = db.db.Exec("UPDATE access SET scope_news = ? WHERE user_id = ?", time.Now(), userId)
		case SCOPE_ZEROTIER:
			_, err = db.db.Exec("UPDATE access SET scope_zerotier = ? WHERE user_id = ?", time.Now(), userId)
		}

		if err != nil {
			return AccessRecord{}, err
		}

		// Get the access key
		err = db.db.QueryRow("SELECT access_key, user_name, scope_news, scope_zerotier, avatar_url FROM access WHERE user_id = ?", userId).Scan(&access.AccessKey, &access.UserName, &access.ScopeNews, &access.ScopeZerotier, &access.UserAvatarUrl)
		if err != nil {
			return AccessRecord{}, err
		}

	} else {
		// Generate a new access key and check if it already exists
		access.AccessKey = randomPassword(userId, userName)
		for db.db.QueryRow("SELECT COUNT(*) FROM access WHERE access_key = ?", userId).Scan(&count); count > 0; {
			access.AccessKey = randomPassword(userId, userName)
		}

		// Insert the access key and set the scope
		switch scope {
		case SCOPE_NEWS:
			_, err = db.db.Exec("INSERT INTO access (user_id, user_name, access_key, avatar_url, scope_news) VALUES (?, ?, ?, ?, ?)", userId, userName, access.AccessKey, userAvatarUrl, time.Now())
			access.ScopeNews = sql.NullTime{Time: time.Now(), Valid: true}
		case SCOPE_ZEROTIER:
			_, err = db.db.Exec("INSERT INTO access (user_id, user_name, access_key, avatar_url, scope_zerotier) VALUES (?, ?, ?, ?, ?)", userId, userName, access.AccessKey, userAvatarUrl, time.Now())
			access.ScopeZerotier = sql.NullTime{Time: time.Now(), Valid: true}
		}

		if err != nil {
			return AccessRecord{}, err
		}
	}

	return access, nil
}

func (db *BotDatabase) AddNetwork(networkId, networkName string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	_, err := db.db.Exec("INSERT INTO network (network_id, network_name) VALUES (?, ?)", networkId, networkName)
	return err
}

func (db *BotDatabase) GetNetwork(networkId string) (ZeroTierNetworkRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var network = ZeroTierNetworkRecord{
		NetworkID: networkId,
	}
	err := db.db.QueryRow("SELECT id, network_name FROM network WHERE network_id = ?", networkId).Scan(&network.ID, &network.NetworkName)
	return network, err
}

func (db *BotDatabase) DeleteNetwork(networkId string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	_, err := db.db.Exec("DELETE FROM network WHERE network_id = ?", networkId)
	return err
}

func (db *BotDatabase) AddMember(networkId, memberId, memberName, user_id string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, err := db.GetMember(networkId, memberId); err == nil {
		return nil
	}

	_, err := db.db.Exec("INSERT INTO member (network_id, member_id, member_name, user_id) VALUES (?, ?, ?, ?)", networkId, memberId, memberName, user_id)
	return err
}

func (db *BotDatabase) GetMember(networkId, memberId string) (ZeroTierMemberRecord, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	var member = ZeroTierMemberRecord{
		NetworkID: networkId,
		MemberID:  memberId,
	}
	err := db.db.QueryRow("SELECT id, member_name, user_id FROM member WHERE network_id = ? AND member_id = ?", networkId, memberId).Scan(&member.ID, &member.MemberName, &member.UserID)
	return member, err
}

func randomPassword(user_id, user_name string) string {

	// Generate a random byte slice
	b := make([]byte, 32)
	rand.Read(b)

	// Generate a random UUID
	u := uuid.New()

	// Concatenate all the data
	b = append(b, []byte(u.String())...)
	b = append(b, []byte(user_id)...)
	b = append(b, []byte(user_name)...)

	// Calculate the SHA1 hash of the message
	h := sha1.New()
	h.Write(b)
	hash := h.Sum(nil)

	// Encode the hash to url safe base64
	return base64.URLEncoding.EncodeToString(hash)
}
