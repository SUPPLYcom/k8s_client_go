package k8s_client_go

import (
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
)

type AksKeyCreds struct {
	ApiVersion string `yaml:"apiVersion"`
	CurrentContext string `yaml:"current-context"`
	Kind string `yaml:"kind"`
	Clusters []struct{
		Name string
		Cluster struct {
			ClientCertificateAuthority string `yaml:"certificate-authority-data"`
			Server string `yaml:"server"`
		}
	}
	Contexts []struct{
		Name string
		Context struct {
			Cluster string `yaml:"cluster"`
			User string `yaml:"user"`
		}
	}
	Users []struct {
		Name string
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
			Token                 string
		}
	}
}

type AksTokenCreds struct {
	ApiVersion string `yaml:"apiVersion"`
	CurrentContext string `yaml:"current-context"`
	Kind string `yaml:"kind"`
	Clusters []struct{
		Name string
		Cluster struct {
			ClientCertificateAuthority string `yaml:"certificate-authority-data"`
			Server string `yaml:"server"`
		}
	}
	Contexts []struct{
		Name string
		Context struct {
			Cluster string `yaml:"cluster"`
			User string `yaml:"user"`
		}
	}
	Users []struct {
		Name string
		User struct {
			AuthProvider struct  {
				Config struct {
					AccessToken string `yaml:"access-token"`
					ApiServerId string `yaml:"apiserver-id"`
					ClientId string `yaml:"client-id"`
					ConfigMode string `yaml:"config-mode"`
					Environment string
					ExpiresIn string `yaml:"expires-in"`
					ExpiresOn string `yaml:"expires-on"`
					RefreshToken string `yaml:"refresh-token"`
					TenantId string `yaml:"tenant-id"`
				} `yaml:"auth-provider"`
			}
		}
	}
}
