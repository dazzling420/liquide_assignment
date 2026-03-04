package redis

import (
	"context"
	"time"

	"github.com/redis/rueidis"
)

func (s RedisRepository) AddSession(key, session, sessionToken string) error {
	err := s.db.Do(context.Background(), s.db.B().Hset().Key(key).FieldValue().FieldValue(session, sessionToken).Build()).Error()
	return err
}

func (s RedisRepository) CheckIfKeyExists(key string) (bool, error) {
	return s.db.Do(context.Background(), s.db.B().Exists().Key(key).Build()).AsBool()
}

func (s RedisRepository) CheckRateLimited(key string) (int64, error) {
	// return value of -2 means key does not exist
	// return value of -1 means key exists but has no ttl
	// return value of > 0 means key exists and has ttl
	return s.db.Do(context.Background(), s.db.B().Ttl().Key(key).Build()).AsInt64()
}

func (s RedisRepository) AddRateLimiting(key string, limit int) {
	var ctx = context.Background()
	s.db.Do(ctx, s.db.B().Set().Key(key).Value("true").Ex(time.Duration(limit)*time.Second).Build())
}

func (s RedisRepository) GetKey(key string) (string, error) {
	resp := s.db.Do(context.Background(), s.db.B().Get().Key(key).Build())
	if resp.Error() != nil {
		if resp.Error() == rueidis.Nil {
			return "Key doesn't exist", rueidis.Nil
		} else {
			return "Internal Server Error", resp.Error()
		}
	}
	r, _ := resp.ToString()
	return r, nil
}

func (s RedisRepository) HGetAll(key string) (map[string]string, error) {
	resp := s.db.Do(context.Background(), s.db.B().Hgetall().Key(key).Build())
	if resp.Error() != nil {
		if resp.Error() == rueidis.Nil {
			return nil, rueidis.Nil
		} else {
			return nil, resp.Error()
		}
	}
	r, _ := resp.AsStrMap()
	return r, nil
}

func (s RedisRepository) GetActiveSessions(key string) (map[string]string, error) {
	p := s.db.Do(context.Background(), s.db.B().Hgetall().Key(key).Build())
	if p.Error() != nil {
		if p.Error() != rueidis.Nil {
			return nil, p.Error()
		} else if p.Error() == rueidis.Nil {
			return nil, rueidis.Nil
		}
	}

	v, err := p.AsStrMap()
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (s RedisRepository) DeleteSession(key, session string) error {
	p := s.db.Do(context.Background(), s.db.B().Hdel().Key(key).Field(session).Build())
	if p.Error() != nil {
		return p.Error()
	}
	return nil
}
