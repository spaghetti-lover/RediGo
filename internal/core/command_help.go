package core

func cmdHELP() []byte {
	helpCommands := []string{
		"--------------------------------",
		"GET key - Get the value of a key",
		"SET key value - Set the value of a key",
		"DEL key - Delete a key",
		"EXISTS key - Check if a key exists",
		"TTL key - Get the time to live for a key",
		"SADD key member [member ...] - Add members to a set",
		"SREM key member [member ...] - Remove members from a set",
		"EXPIRE key seconds - Set a timeout on a key",
		"HELP - Show this help message",
		"--------------------------------",
	}
	return Encode(helpCommands, false)
}
