package cache

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v7"
	"strings"
	"xlink/user_service/internal/utils"
)

const (
	userIdKeyPrefix = "user-id-"
	tokenKeyPrefix  = "token-"
	roleKeyPrefix   = "role-"
	sep             = ";"
)

type UserCacheRepositoryRedis struct {
	RedisClient *redis.Client
}

func getUserIdKey(userId string) string {
	return userIdKeyPrefix + userId
}

func getTokenKey(token string) string {
	return tokenKeyPrefix + token
}

func getRolesKey(userId string) string {
	return roleKeyPrefix + userId
}

func boolToString(val bool) string {
	switch val {
	case true:
		return "1"
	default:
		return "0"
	}
}

func stringToBool(val string) bool {
	switch val {
	case "1":
		return true
	default:
		return false
	}
}

func NewUserCacheRepositoryRedis(redisClient *redis.Client) *UserCacheRepositoryRedis {
	return &UserCacheRepositoryRedis{RedisClient: redisClient}
}

func (t *UserCacheRepositoryRedis) CheckToken(userId string, token string) (bool, error) {
	commandResult := t.RedisClient.Get(getUserIdKey(userId))
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return false, nil
		}
		return false, commandResult.Err()
	}
	if commandResult.Val() != token {
		return false, nil
	}
	return true, nil
}

func (t *UserCacheRepositoryRedis) GetToken(userId string) (string, error) {
	commandResult := t.RedisClient.Get(getUserIdKey(userId))
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return "", nil
		}
		return "", commandResult.Err()
	}
	return commandResult.Val(), nil
}

func (t *UserCacheRepositoryRedis) SetToken(userId string, token string) error {
	commandResult := t.RedisClient.Set(getUserIdKey(userId), token, 0)
	if commandResult.Err() != nil {
		return fmt.Errorf("couldn't cache token in redis: %v", commandResult.Err())
	}

	return nil
}

func (t *UserCacheRepositoryRedis) GetRole(userId string) (string, bool, bool, error) {
	commandResult := t.RedisClient.Get(getRolesKey(userId))
	if commandResult.Err() != nil {
		return "", false, false, fmt.Errorf("couldn't get roles from redis: %v", commandResult.Err())
	}

	rolesStrings := strings.Split(commandResult.Val(), sep)
	if len(rolesStrings) != 2 {
		return "", false, false,
			fmt.Errorf("couldn't get roles from redis: invalid format (expected 2 strings divided by '%s', got '%s'",
				sep, commandResult.Val())
	}
	isStaffString := rolesStrings[0]
	isAdminString := rolesStrings[1]

	isStaff := stringToBool(isStaffString)
	isAdmin := stringToBool(isAdminString)

	role := utils.GetRoleByIsStaffIsAdmin(isStaff, isAdmin)

	return role, isStaff, isAdmin, nil
}

func (t *UserCacheRepositoryRedis) SetRole(userId string, isStaff bool, isAdmin bool) error {
	rolesString := fmt.Sprintf("%s%s%s", boolToString(isStaff), sep, boolToString(isAdmin))

	commandResult := t.RedisClient.Set(getRolesKey(userId), rolesString, 0)
	if commandResult.Err() != nil {
		return fmt.Errorf("couldn't cache roles in redis: %v", commandResult.Err())
	}
	return nil
}
