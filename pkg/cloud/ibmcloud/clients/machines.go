/*
Copyright 2018 The Kubernetes authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package clients

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

type GuestService struct {
	sess *session.Session
}

func NewGuestService(sess *session.Session) GuestService {
	return GuestService{sess: sess}
}

type IBMCloudConfig struct {
	UserName string `yaml:"userName,omitempty"`
	APIKey   string `yaml:"apiKey,omitempty"`
}

func NewInstanceServiceFromMachine(kubeClient kubernetes.Interface, machine *clusterv1.Machine) (GuestService, error) {
	// clouds.yaml is mounted into controller pod for clouds authentication
	fileName := "/etc/ibmcloud/clouds.yaml"
	if _, err := os.Stat(fileName); err != nil {
		return GuestService{}, fmt.Errorf("Cannot stat %q: %v", fileName, err)
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return GuestService{}, fmt.Errorf("Cannot read %q: %v", fileName, err)
	}

	config := IBMCloudConfig{}
	yaml.Unmarshal(bytes, &config)

	if config.UserName == "" || config.APIKey == "" {
		return GuestService{}, fmt.Errorf("Failed getting IBM Cloud config userName %q, apiKey %q", config.UserName, config.APIKey)
	}

	sess := session.New(config.UserName, config.APIKey)
	return NewGuestService(sess), nil
}

func (gs *GuestService) guestWaitReady(Id int) {
	// Wait for transactions to finish
	log.Printf("Waiting for transactions to complete before destroying.")
	s := services.GetVirtualGuestService(gs.sess).Id(Id)

	// Delay to allow transactions to be registered
	time.Sleep(5 * time.Second)

	for transactions, _ := s.GetActiveTransactions(); len(transactions) > 0; {
		log.Print(".")
		// TODO(gyliu513) make it configurable or use the notification mechanism to optimize
		// the process instead of hardcoded waiting.
		time.Sleep(5 * time.Second)
		transactions, _ = s.GetActiveTransactions()
	}
	log.Println("wait done")
}

func (gs *GuestService) GuestCreate(clusterName, hostName, sshKeyName, userScript string) {
	s := services.GetVirtualGuestService(gs.sess)

	// TODO: use customized value instead of hardcoded in code
	keyId := getSshKey(gs.sess, sshKeyName)
	if keyId == 0 {
		log.Printf("Cannot retrieving specific SSH key. Stop creating VM instance\n")
		return
	}
	sshKeys := []datatypes.Security_Ssh_Key{
		{
			Id: sl.Int(keyId),
		},
	}

	// TODO: remove this hardcoded path if we know how to execute userData as a script directly
	scriptURL := "https://raw.githubusercontent.com/multicloudlab/cluster-api-provider-ibmcloud/master/cmd/clusterctl/configuration/ibmcloud/scripts/launcher.sh"
	userData := []datatypes.Virtual_Guest_Attribute{
		{
			// TODO: if base64 needed
			Value: sl.String(userScript),
			Guest: nil,
			Type: &datatypes.Virtual_Guest_Attribute_Type{
				Keyname: sl.String("USER_DATA"),
				Name:    sl.String("user data"),
			},
		},
	}

	// Create a Virtual_Guest instance as a template
	vGuestTemplate := datatypes.Virtual_Guest{
		Hostname:                     sl.String(hostName),
		Domain:                       sl.String("example.com"),
		MaxMemory:                    sl.Int(4096),
		StartCpus:                    sl.Int(1),
		Datacenter:                   &datatypes.Location{Name: sl.String("wdc01")},
		OperatingSystemReferenceCode: sl.String("UBUNTU_LATEST"),
		LocalDiskFlag:                sl.Bool(true),
		HourlyBillingFlag:            sl.Bool(true),
		SshKeyCount:                  sl.Uint(1),
		SshKeys:                      sshKeys,
		PostInstallScriptUri:         sl.String(scriptURL),
		UserData:                     userData,
	}

	vGuest, err := s.Mask("id;domain").CreateObject(&vGuestTemplate)
	if err != nil {
		log.Printf("%s\n", err)
		return
	} else {
		log.Printf("\nNew Virtual Guest created with ID %d\n", *vGuest.Id)
		log.Printf("Domain: %s\n", *vGuest.Domain)
	}

	// Wait for transactions to finish
	log.Printf("Waiting for transactions to complete before destroying.")
	gs.guestWaitReady(*vGuest.Id)
}

func (gs *GuestService) GuestDelete(Id int) error {
	s := services.GetVirtualGuestService(gs.sess).Id(Id)

	success, err := s.DeleteObject()
	if err != nil {
		log.Printf("Error deleting virtual guest: %s", err)
	} else if success == false {
		log.Printf("Error deleting virtual guest")
		err = fmt.Errorf("Error delete virtual guest")
	} else {
		log.Printf("Virtual Guest deleted successfully")
	}
	return err
}

func (gs *GuestService) GuestList() ([]datatypes.Virtual_Guest, error) {
	s := services.GetAccountService(gs.sess)

	guests, err := s.GetVirtualGuests()
	if err != nil {
		log.Printf("Error listinging virtual guest: %s", err)
		return []datatypes.Virtual_Guest{}, err
	}
	return guests, nil
}

// FIXME: use API layer query instead of query all then compare here
func (gs *GuestService) GuestGet(name string) (*datatypes.Virtual_Guest, error) {
	var vg *datatypes.Virtual_Guest
	guests, err := gs.GuestList()

	if err != nil {
		return vg, err
	}

	for _, guest := range guests {
		// FIXME: how to unique identify one guest
		if *guest.Hostname == name {
			log.Printf("Found guest with Id %d for %s", *guest.Id, name)
			return &guest, nil
		}
	}
	return vg, fmt.Errorf("unable to find guest with name %s", name)
}

func getSshKey(sess *session.Session, name string) int {
	id := 0

	service := services.GetAccountService(sess)
	keys, err := service.GetSshKeys()
	if err != nil {
		log.Printf("Error retrieving ssh keys from Account: %s\n", err)
		return id
	}

	for _, key := range keys {
		if *key.Label == name {
			id = *key.Id
			log.Printf("Get SSH key for %q with value %d\n", *key.Label, *key.Id)
			break
		}
	}

	return id

}