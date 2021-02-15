package k8s_client_go

import (
	"encoding/base64"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// HostURL - Default Hashicups URL
const HostURL string = "http://localhost:19090"

// Client -
type Client struct {
	HostURL    string
	K8sClientSet *kubernetes.Clientset
}

func NewClientFromKubeFile(apihost, kubeconfig string) (*Client, error) {
	config, configErr := clientcmd.BuildConfigFromFlags(apihost, kubeconfig)
	if(configErr != nil) {
		return nil, configErr
	}
	clientset, clientsetErr := kubernetes.NewForConfig(config)
	if(clientsetErr != nil) {
		return nil, clientsetErr
	}

	newclient := Client{
		HostURL: apihost,
		K8sClientSet: clientset,
	}

	return &newclient, nil
}

func NewClientFromKubeCreds(clusterName, contextName, userName, apiHost, clientCertificateAuthority, clientCertificateData, clientKeyData, token string) (*Client, error) {
	tmpKubeConf, _ := GenerateKubeFileFromKey(clusterName, contextName, userName, apiHost, clientCertificateAuthority, clientCertificateData, clientKeyData, token)
	defer os.Remove(tmpKubeConf.Name())

	config, configErr := clientcmd.BuildConfigFromFlags(apiHost, tmpKubeConf.Name())
	if(configErr != nil) {
		return nil, configErr
	}
	clientset, clientsetErr := kubernetes.NewForConfig(config)
	if(clientsetErr != nil) {
		return nil, clientsetErr
	}

	newclient := Client{
		HostURL:      apiHost,
		K8sClientSet: clientset,
	}

	return &newclient, nil
}

func NewClientFromToken(host, certAuthority, bearToken string) (*Client, error) {
	certData, certDataErr := base64.StdEncoding.DecodeString(certAuthority)
	if(certData == nil || certDataErr != nil) {
		return &Client{}, errors.New("error decoding cert auth data")
	}

	config := rest.Config{
		Host:                host,
		BearerToken:         bearToken,
		TLSClientConfig:     rest.TLSClientConfig{
			Insecure:   false,
			CAData:     certData,
		},
	}

	clientset, clientSetErr := kubernetes.NewForConfig(&config)
	if clientset == nil || clientSetErr != nil {
		return &Client{}, errors.New("error initializing client")
	}

	newclient := Client{
		HostURL: host,
		K8sClientSet: clientset,
	}

	return &newclient, nil
}

func GenerateKubeConf(name, apiHost, clientCertificateAuthority, clientCertificateData, clientKeyData, token  string) (string, error) {
	aksCreds := AksKeyCreds{
		ApiVersion:     "v1",
		CurrentContext: "tfprovtest",
		Kind:           "Config",
		Clusters: []struct {
			Name    string
			Cluster struct {
				ClientCertificateAuthority string `yaml:"certificate-authority-data"`
				Server string `yaml:"server"`
			}
		}{
			{
				"tfprovtest",
				struct {
					ClientCertificateAuthority string `yaml:"certificate-authority-data"`
					Server string `yaml:"server"`
				}{
					ClientCertificateAuthority: clientCertificateAuthority,
					Server: apiHost,
				},
			},
		},
		Contexts: []struct {
			Name    string
			Context struct {
				Cluster string `yaml:"cluster"`
				User    string `yaml:"user"`
			}
		}{
			{
				"tfprovtest",
				struct {
					Cluster string `yaml:"cluster"`
					User string `yaml:"user"`
				}{
					"tfprovtest",
					"clusterUser_feiseu2-supply-rg-001_feiseu2supplyaks",
				},
			},
		},
		Users: []struct {
			Name string
			User struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKeyData         string `yaml:"client-key-data"`
				Token                 string
			}
		}{
			{
				name,
				struct {
					ClientCertificateData string `yaml:"client-certificate-data"`
					ClientKeyData         string `yaml:"client-key-data"`
					Token                 string
				}{
					ClientCertificateData: clientCertificateData,
					ClientKeyData:         clientKeyData,
					Token:                 token,
				},
			},
		},
	}

	aksCredsYaml, _ := yaml.Marshal(aksCreds)
	aksCredsYamlStr := string(aksCredsYaml)

	return aksCredsYamlStr, nil
}

func GenerateKubeFileFromKey(clusterName, contextName, userName, apiHost, clientCertificateAuthority, clientCertificateData, clientKeyData, token string) (*os.File, error) {
	aksCreds := AksKeyCreds{
		ApiVersion:     "v1",
		CurrentContext: "tfprovtest",
		Kind:           "Config",
		Clusters: []struct {
			Name    string
			Cluster struct {
				ClientCertificateAuthority string `yaml:"certificate-authority-data"`
				Server string `yaml:"server"`
			}
		}{
			{
				clusterName,
				struct {
					ClientCertificateAuthority string `yaml:"certificate-authority-data"`
					Server string `yaml:"server"`
				}{
					ClientCertificateAuthority: clientCertificateAuthority,
					Server: apiHost,
				},
			},
		},
		Contexts: []struct {
			Name    string
			Context struct {
				Cluster string `yaml:"cluster"`
				User    string `yaml:"user"`
			}
		}{
			{
				contextName,
				struct {
					Cluster string `yaml:"cluster"`
					User string `yaml:"user"`
				}{
					clusterName,
					userName,
				},
			},
		},
		Users: []struct {
			Name string
			User struct {
				ClientCertificateData string `yaml:"client-certificate-data"`
				ClientKeyData         string `yaml:"client-key-data"`
				Token                 string
			}
		}{
			{
				userName,
				struct {
					ClientCertificateData string `yaml:"client-certificate-data"`
					ClientKeyData         string `yaml:"client-key-data"`
					Token                 string
				}{
					ClientCertificateData: clientCertificateData,
					ClientKeyData:         clientKeyData,
					Token:                 token,
				},
			},
		},
	}

	aksCredsYaml, _ := yaml.Marshal(aksCreds)
	println(string(aksCredsYaml))

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return nil, errors.New("Cannot create temporary file")
	}

	if _, err = tmpFile.Write(aksCredsYaml); err != nil {
		return nil, errors.New("Failed to write to temporary file")
	}

	return tmpFile, nil
}

func GenerateKubeFileFromToken(clusterName, contextName, userName, apiHost, clientCertificateAuthority, accessToken, apiServerId, clientId, configMode, environment, expiresIn, expiresOn, refreshToken, tenantId string) (*os.File, error) {
	aksCreds := AksTokenCreds{
		ApiVersion:     "v1",
		CurrentContext: "tfprovtest",
		Kind:           "Config",
		Clusters: []struct {
			Name    string
			Cluster struct {
				ClientCertificateAuthority string `yaml:"certificate-authority-data"`
				Server                     string `yaml:"server"`
			}
		}{
			{
				clusterName,
				struct {
					ClientCertificateAuthority string `yaml:"certificate-authority-data"`
					Server string `yaml:"server"`
				}{
					ClientCertificateAuthority: clientCertificateAuthority,
					Server: apiHost,
				},
			},
		},
		Contexts: []struct {
			Name    string
			Context struct {
				Cluster string `yaml:"cluster"`
				User    string `yaml:"user"`
			}
		}{
			{
				contextName,
				struct {
					Cluster string `yaml:"cluster"`
					User string `yaml:"user"`
				}{
					clusterName,
					userName,
				},
			},
		},
		Users: []struct {
			Name string
			User struct {
				AuthProvider struct {
					Config struct {
						AccessToken  string `yaml:"access-token"`
						ApiServerId  string `yaml:"apiserver-id"`
						ClientId     string `yaml:"client-id"`
						ConfigMode   string `yaml:"config-mode"`
						Environment  string
						ExpiresIn    string `yaml:"expires-in"`
						ExpiresOn    string `yaml:"expires-on"`
						RefreshToken string `yaml:"refresh-token"`
						TenantId     string `yaml:"tenant-id"`
					} `yaml:"auth-provider"`
				}
			}
		}{
			{
				Name: userName,
				User: struct {
					AuthProvider struct {
						Config struct {
							AccessToken  string `yaml:"access-token"`
							ApiServerId  string `yaml:"apiserver-id"`
							ClientId     string `yaml:"client-id"`
							ConfigMode   string `yaml:"config-mode"`
							Environment  string
							ExpiresIn    string `yaml:"expires-in"`
							ExpiresOn    string `yaml:"expires-on"`
							RefreshToken string `yaml:"refresh-token"`
							TenantId     string `yaml:"tenant-id"`
						} `yaml:"auth-provider"`
					}
				}{
					struct {
						Config struct {
							AccessToken  string `yaml:"access-token"`
							ApiServerId  string `yaml:"apiserver-id"`
							ClientId     string `yaml:"client-id"`
							ConfigMode   string `yaml:"config-mode"`
							Environment  string
							ExpiresIn    string `yaml:"expires-in"`
							ExpiresOn    string `yaml:"expires-on"`
							RefreshToken string `yaml:"refresh-token"`
							TenantId     string `yaml:"tenant-id"`
						}`yaml:"auth-provider"`
					}{
						struct {
							AccessToken string `yaml:"access-token"`
							ApiServerId string `yaml:"apiserver-id"`
							ClientId string `yaml:"client-id"`
							ConfigMode string `yaml:"config-mode"`
							Environment string
							ExpiresIn string `yaml:"expires-in"`
							ExpiresOn string `yaml:"expires-on"`
							RefreshToken string `yaml:"refresh-token"`
							TenantId string `yaml:"tenant-id"`
						} {
							accessToken,
							apiServerId,
							clientId,
							configMode,
							environment,
							expiresIn,
							expiresOn,
							refreshToken,
							tenantId,
						},
					},
				},
			},
		},
	}

	aksCredsYaml, _ := yaml.Marshal(aksCreds)
	println(string(aksCredsYaml))

	tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-")
	if err != nil {
		return nil, errors.New("Cannot create temporary file")
	}

	if _, err = tmpFile.Write(aksCredsYaml); err != nil {
		return nil, errors.New("Failed to write to temporary file")
	}

	return tmpFile, nil
}