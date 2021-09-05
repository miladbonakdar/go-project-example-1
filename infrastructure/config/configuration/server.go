package configuration

type ServerConfiguration struct {
	Port                   int    // this is for server port inside the container
	OutSideOfContainerPort int    // this is for the port that is observable from outside the container. for swagger gen
	OutSideOfContainerHost string // this is the hostname outside of the container. for swagger gen. default is localhost
}
