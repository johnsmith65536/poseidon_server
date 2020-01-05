package redis

const HeartBeatChannel = "poseidon_heart_beat_channel"

func HeartBeat() error {
	return redisCli.Publish(HeartBeatChannel, "").Err()
}
