package apikeydomain

func (a APIKey) ID() APIKeyID                   { return a.id }
func (a APIKey) OwnerID() OwnerID               { return a.ownerID }
func (a APIKey) KeyName() APIKeyName            { return a.name }
func (a APIKey) Provider() ProviderName         { return a.provider }
func (a APIKey) ModelName() ModelName           { return a.model }
func (a APIKey) UsageEnv() UsageEnv             { return a.usageEnv }
func (a APIKey) KeyValueEnc() EncryptedKeyValue { return a.keyValueEnc }
