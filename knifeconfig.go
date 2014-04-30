package chef

import "crypto/rsa"

// http://docs.opscode.com/config_rb_knife.html
type KnifeConfig struct {
	ChefServerUrl string
	Host          string
	Port          string
	//chef_zero[:enabled]
	//chef_zero[:port]
	ClientKey         *rsa.PrivateKey
	CookbookCopyright string
	CookbookEmail     string
	CookbookLicense   string
	//cookbook_path map[]string
	DataBagEncryptVersion int64
	LocalMode             bool
	NodeName              string
	//no_proxy map[]string
	SyntaxCheckCachePath string
	ValidationClientName string
	ValidationKey        string
	VersionedCookbooks   bool
}
