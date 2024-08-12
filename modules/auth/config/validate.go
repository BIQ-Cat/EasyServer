package config

func init() {
	if Config.Create.RequireData == CREATE_USERNAME_ONLY_REQUIRED && !Config.Create.HasUsername {
		panic(`You must enable username if nither email nor phone number required.
			  -> Fix: set Config.Create.HasUsername to true or choose another constant`)
	}
}
