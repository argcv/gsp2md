package config

type SMTPConfig = struct {
	Host   string
	Port   int
	Sender string
	User   string
	Pass   string
	// SkipTLS
	InsecureSkipVerify bool
}

/*
Example Config File:

	mail:
		smtp:
			default: # profile
				host: 127.0.0.1
				port: 25
				sender: Alice Ez
				user: "aez@example.com"
				pass: "your-password" # empty for no authorization
				option:
					insecure_skip_verify: true
 */
func LoadSMTPConfig(profile string) *SMTPConfig {
	ckb := NewKeyBuilder(KeyMailBase, KeyMailSMTPBase).WithProfile(profile)
	opt := ckb.Clone().WithClass(KeyMailSMTPOption)
	return &SMTPConfig{
		Host:               ckb.GetStringOrDefault(KeyMailSMTPHost, "127.0.0.1"),
		Port:               ckb.GetIntOrDefault(KeyMailSMTPPort, 25),
		Sender:             ckb.GetStringOrDefault(KeyMailSMTPSender, ""),
		User:               ckb.GetStringOrDefault(KeyMailSMTPUser, "aez@example.com"),
		Pass:               ckb.GetStringOrDefault(KeyMailSMTPPass, "pass"),
		InsecureSkipVerify: opt.GetBoolOrDefault(KeyMailSMTPOptionInsecureSkipVerify, false),
	}
}
